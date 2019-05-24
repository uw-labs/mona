package command

import (
	"github.com/davidsbond/mona/internal/files"
	"github.com/davidsbond/mona/internal/hash"
	"os/exec"
	"strings"
)

// Build will execute the build commands for all modules where changes
// are detected.
func Build() error {
	changed, err := getChangedModules()

	if err != nil {
		return err
	}

	newHashes := make(map[string]string)
	for _, module := range changed {
		parts := strings.Split(module.Commands.Build, " ")
		cmd := exec.Command(parts[0], parts[1:]...)
		cmd.Dir = module.Location

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

		newHash, err := hash.Generate(module.Location, module.Exclude...)

		if err != nil {
			return err
		}

		newHashes[module.Name] = newHash
	}

	lock, err := files.LoadLockFile()

	if err != nil {
		return err
	}

	for i, lockInfo := range lock.Modules {
		name, location, _ := files.ParseLockLine(lockInfo)

		if hash, ok := newHashes[name]; ok {
			lock.Modules[i] = files.CreateLockLine(name, location, hash)
		}
	}

	return files.UpdateLockFile(lock)
}
