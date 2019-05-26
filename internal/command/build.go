package command

import (
	"github.com/davidsbond/mona/internal/files"
	"github.com/davidsbond/mona/internal/hash"
)

// Build will execute the build commands for all modules where changes
// are detected.
func Build() error {
	project, err := files.LoadProjectFile()

	if err != nil {
		return err
	}

	changed, err := getChangedModules(changeTypeBuild)

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

		newHashes[module.Name+module.Location] = newHash

		if len(module.Artefacts) > 0 {
			if err := module.CollectArtefacts(project.Artefacts); err != nil {
				return err
			}
		}

	}

	lock, err := files.LoadLockFile()

	if err != nil {
		return err
	}

	for i, lockInfo := range lock.Modules {
		if hash, ok := newHashes[lockInfo.Name+lockInfo.Location]; ok {
			lock.Modules[i].BuildHash = hash
		}
	}

	return files.UpdateLockFile(lock)
}
