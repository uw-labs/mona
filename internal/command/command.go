// Package command contains methods called by the CLI to manage
// a mona project.
package command

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
	"strings"

	"github.com/davidsbond/mona/internal/config"

	"github.com/davidsbond/mona/internal/files"
	"github.com/davidsbond/mona/pkg/hashdir"
)

type (
	changeType int
	rangeFn    func(*files.ModuleFile) error
)

const (
	changeTypeBuild changeType = 0
	changeTypeTest  changeType = 1
	changeTypeLint  changeType = 2
)

func getChangedModules(pj *files.ProjectFile, change changeType) ([]*files.ModuleFile, error) {
	modules, err := files.FindModules(pj.Location)

	if err != nil {
		return nil, err
	}

	lock, err := files.LoadLockFile(pj.Location)

	if err != nil {
		return nil, err
	}

	var out []*files.ModuleFile
	for _, modInfo := range modules {
		lockInfo, ok := lock.Modules[modInfo.Name]

		if !ok {
			// If the module isn't present in the lock file, it's probably
			// new and needs adding to the lock file.
			out = append(out, modInfo)
			continue
		}

		// Generate a new hash for the module directory
		exclude := append(pj.Exclude, modInfo.Exclude...)
		newHash, err := hashdir.Generate(modInfo.Location, config.Parallelism, exclude...)

		if err != nil {
			return nil, err
		}

		// Check if we need to build/lint/test
		diff := false
		switch change {
		case changeTypeBuild:
			diff = lockInfo.BuildHash != newHash
		case changeTypeTest:
			diff = lockInfo.TestHash != newHash
		case changeTypeLint:
			diff = lockInfo.LintHash != newHash
		}

		if diff {
			out = append(out, modInfo)
		}
	}

	return out, nil
}

func streamOutputs(outputs ...io.Reader) {
	for _, output := range outputs {
		go func(o io.Reader) {
			scanner := bufio.NewScanner(o)
			scanner.Split(bufio.ScanLines)

			for scanner.Scan() {
				m := scanner.Text()
				fmt.Println(m)
			}
		}(output)
	}
}

func streamCommand(command, wd string) error {
	parts := strings.Split(command, " ")
	cmd := exec.Command(parts[0], parts[1:]...)
	cmd.Dir = wd

	stdout, err := cmd.StdoutPipe()

	if err != nil {
		return err
	}

	stderr, err := cmd.StderrPipe()

	if err != nil {
		return err
	}

	streamOutputs(stdout, stderr)

	if err := cmd.Start(); err != nil {
		return err
	}

	return cmd.Wait()
}

func rangeChangedModules(pj *files.ProjectFile, change changeType, fn rangeFn, updateHashes bool) error {
	changed, err := getChangedModules(pj, change)

	if err != nil || len(changed) == 0 {
		return err
	}

	lock, err := files.LoadLockFile(pj.Location)

	if err != nil {
		return err
	}

	for _, module := range changed {
		if err := fn(module); err != nil {
			return err
		}

		if !updateHashes {
			continue
		}

		exclude := append(pj.Exclude, module.Exclude...)
		newHash, err := hashdir.Generate(module.Location, config.Parallelism, exclude...)

		if err != nil {
			return err
		}

		lockInfo, modInLock := lock.Modules[module.Name]

		if !modInLock {
			if err := files.AddModule(lock, pj.Location, module.Name); err != nil {
				return err
			}

			lockInfo = lock.Modules[module.Name]
		}

		switch change {
		case changeTypeBuild:
			lockInfo.BuildHash = newHash
		case changeTypeTest:
			lockInfo.TestHash = newHash
		case changeTypeLint:
			lockInfo.LintHash = newHash
		}

	}

	return files.UpdateLockFile(pj.Location, lock)
}
