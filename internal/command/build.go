package command

import (
	"github.com/apex/log"

	"github.com/davidsbond/mona/internal/files"
)

// Build will execute the build commands for all modules where changes
// are detected.
func Build(pj *files.ProjectFile) error {
	return rangeChangedModules(pj, buildModule, changeTypeBuild)
}

func buildModule(module *files.ModuleFile) error {
	log.Infof("Building module %s at %s", module.Name, module.Location)

	if module.Commands.Build == "" {
		return nil
	}

	return streamCommand(module.Commands.Build, module.Location)
}
