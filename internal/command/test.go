package command

import (
	"github.com/davidsbond/mona/internal/files"
	"github.com/davidsbond/mona/internal/hash"
)

// Test attempts to run the test command for all modules where changes
// are detected.
func Test() error {
	changed, err := getChangedModules(changeTypeTest)

	if err != nil {
		return err
	}

	newHashes := make(map[string]string)
	for _, module := range changed {
		if err := testModule(module); err != nil {
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
			lock.Modules[i].TestHash = hash
		}
	}

	return files.UpdateLockFile(lock)
}
