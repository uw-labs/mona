package command

import (
	"github.com/davidsbond/mona/internal/files"
	"github.com/davidsbond/mona/internal/hash"
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
		if err := buildModule(module); err != nil {
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
