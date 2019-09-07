// Package command contains methods called by the CLI to manage
// a mona project.
package command

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
	"strings"

	"github.com/apex/log"
	"github.com/davidsbond/mona/internal/deps"
	"github.com/davidsbond/mona/internal/files"
	"github.com/davidsbond/mona/internal/hash"
	"github.com/hashicorp/go-multierror"
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

func getChangedModules(mod deps.Module, pj *files.ProjectFile, change changeType) ([]*files.ModuleFile, error) {
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

		// GenerateString a new hash for the module directory
		exclude := append(pj.Exclude, modInfo.Exclude...)
		newHash, err := hash.GetForApp(mod, modInfo.Location, exclude...)

		if err != nil {
			return nil, err
		}

		// Check if we need to build/lint/test
		var oldHash string
		switch change {
		case changeTypeBuild:
			oldHash = lockInfo.BuildHash
		case changeTypeTest:
			oldHash = lockInfo.TestHash
		case changeTypeLint:
			oldHash = lockInfo.LintHash
		}

		if oldHash != newHash {
			out = append(out, modInfo)
		}
	}

	return out, nil
}

func streamOutputs(outputs ...io.ReadCloser) {
	for _, output := range outputs {
		go func(o io.ReadCloser) {
			defer o.Close()

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

func rangeChangedModules(mod deps.Module, pj *files.ProjectFile, fn rangeFn, ct changeType) error {
	changed, err := getChangedModules(mod, pj, ct)

	if err != nil || len(changed) == 0 {
		return err
	}

	lock, err := files.LoadLockFile(pj.Location)

	if err != nil {
		return err
	}

	var errs []error
	for _, module := range changed {
		if err := fn(module); err != nil {
			errs = append(errs, fmt.Errorf("module %s: %s", module.Name, err.Error()))
			continue
		}

		exclude := append(pj.Exclude, module.Exclude...)
		newHash, err := hash.GetForApp(mod, module.Location, exclude...)

		if err != nil {
			return err
		}

		lockInfo, modInLock := lock.Modules[module.Name]

		if !modInLock {
			log.Debugf("Detected new module %s at %s, adding to lock file", module.Name, module.Location)

			if err := files.AddModule(lock, pj.Location, module.Name); err != nil {
				return err
			}

			lockInfo = lock.Modules[module.Name]
		}

		switch ct {
		case changeTypeBuild:
			lockInfo.BuildHash = newHash
		case changeTypeTest:
			lockInfo.TestHash = newHash
		case changeTypeLint:
			lockInfo.LintHash = newHash
		}

		if err := files.UpdateLockFile(pj.Location, lock); err != nil {
			return err
		}
	}

	return multierror.Append(nil, errs...).ErrorOrNil()
}
