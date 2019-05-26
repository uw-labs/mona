package command

import (
	"github.com/davidsbond/mona/internal/files"
	"github.com/davidsbond/mona/internal/hash"
)

// Lint iterates over all new/modified modules and executes their lint command. Once complete,
// the lint hash is updated in the lock file.
func Lint() error {
	changed, err := getChangedModules(changeTypeLint)

	if err != nil {
		return err
	}

	newHashes := make(map[string]string)
	for _, module := range changed {
		if err := lintModule(module); err != nil {
			return err
		}

		newHash, err := hash.Generate(module.Location, module.Exclude...)

		if err != nil {
			return err
		}

		newHashes[module.Name+module.Location] = newHash
	}

	lock, err := files.LoadLockFile()

	if err != nil {
		return err
	}

	for i, lockInfo := range lock.Modules {
		if hash, ok := newHashes[lockInfo.Name+lockInfo.Location]; ok {
			lock.Modules[i].LintHash = hash
		}
	}

	return files.UpdateLockFile(lock)
}
