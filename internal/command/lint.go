package command

import (
	"github.com/davidsbond/mona/internal/files"
)

// Lint iterates over all new/modified modules and executes their lint command. Once complete,
// the lint hash is updated in the lock file.
func Lint(wd string) error {
	return rangeChangedModules(wd, changeTypeLint, lintModule, true)
}

func lintModule(wd string, module *files.ModuleFile) error {
	if module.Commands.Lint == "" {
		return nil
	}

	return streamCommand(module.Commands.Lint, module.Location)
}
