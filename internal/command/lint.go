package command

import (
	"github.com/davidsbond/mona/internal/files"
)

// Lint iterates over all new/modified modules and executes their lint command. Once complete,
// the lint hash is updated in the lock file.
func Lint(pj *files.ProjectFile, parallelism int) error {
	return rangeChangedModules(pj, lintModule, rangeOptions{
		changeType:   changeTypeLint,
		updateHashes: true,
		parallelism:  parallelism,
	})
}

func lintModule(module *files.ModuleFile) error {
	if module.Commands.Lint == "" {
		return nil
	}

	return streamCommand(module.Commands.Lint, module.Location)
}
