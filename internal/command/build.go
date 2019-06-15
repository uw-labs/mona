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

	if module.Commands.Build.Run == "" {
		return nil
	}

	// Run command locally if no image is set
	if module.Commands.Build.Image == "" {
		return streamCommand(module.Commands.Build.Run, module.Location)
	}

	return runInImage(
		module.Commands.Build.Image,
		module.Commands.Build.Run,
		module.Location,
	)
}
