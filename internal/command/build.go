package command

import (
	"github.com/davidsbond/mona/internal/files"
)

// Build will execute the build commands for all modules where changes
// are detected.
func Build() error {
	return rangeChangedModules(changeTypeBuild, true, buildModule)
}

func buildModule(module *files.ModuleFile) error {
	if module.Commands.Build == "" {
		return nil
	}

	return streamCommand(module.Commands.Build, module.Location)
}
