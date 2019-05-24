// Package command contains methods called by the CLI to manage
// a mona project.
package command

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
	"strings"

	"github.com/davidsbond/mona/internal/files"
	"github.com/davidsbond/mona/internal/hash"
)

type (
	changeType int
)

const (
	changeTypeBuild changeType = 0
	changeTypeTest  changeType = 1
)

func getChangedModules(change changeType) ([]*files.ModuleFile, error) {
	lock, err := files.LoadLockFile()

	if err != nil {
		return nil, err
	}

	var out []*files.ModuleFile
	for _, lockInfo := range lock.Modules {
		module, err := files.LoadModuleFile(lockInfo.Location)

		if err != nil {
			return nil, err
		}

		newHash, err := hash.Generate(lockInfo.Location, module.Exclude...)

		if err != nil {
			return nil, err
		}

		diff := false
		switch change {
		case changeTypeBuild:
			diff = lockInfo.BuildHash != newHash
		case changeTypeTest:
			diff = lockInfo.TestHash != newHash
		}

		if diff {
			out = append(out, module)
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

func buildModule(module *files.ModuleFile) error {
	parts := strings.Split(module.Commands.Build, " ")
	cmd := exec.Command(parts[0], parts[1:]...)
	cmd.Dir = module.Location

	return streamCommand(cmd)
}

func testModule(module *files.ModuleFile) error {
	parts := strings.Split(module.Commands.Test, " ")
	cmd := exec.Command(parts[0], parts[1:]...)
	cmd.Dir = module.Location

	return streamCommand(cmd)
}

func streamCommand(cmd *exec.Cmd) error {
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

	if err := cmd.Wait(); err != nil {
		return err
	}

	return nil
}
