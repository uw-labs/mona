package command

import (
	"context"

	"github.com/apex/log"
	"github.com/davidsbond/mona/internal/deps"
	"github.com/davidsbond/mona/internal/environment"
	"github.com/davidsbond/mona/internal/files"
)

// Build will execute the build commands for all modules where changes
// are detected.
func Build(mod deps.Module, pj *files.ProjectFile) error {
	return rangeChangedModules(mod, pj, buildModule, changeTypeBuild)
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

	env, err := environment.NewDockerEnvironment(
		module.Commands.Build.Image,
		module.Location,
	)

	if err != nil {
		return err
	}

	ctx := context.Background()
	return env.Execute(ctx, module.Commands.Build.Run)
}
