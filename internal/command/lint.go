package command

import (
	"context"

	"github.com/apex/log"
	"github.com/davidsbond/mona/internal/environment"
	"github.com/davidsbond/mona/internal/files"
)

// Lint iterates over all new/modified modules and executes their lint command. Once complete,
// the lint hash is updated in the lock file.
func Lint(pj *files.ProjectFile) error {
	return rangeChangedModules(pj, lintModule, changeTypeLint)
}

func lintModule(module *files.ModuleFile) error {
	log.Infof("Linting module %s at %s", module.Name, module.Location)

	if module.Commands.Lint.Run == "" {
		return nil
	}

	// Run command locally if no image is specified
	if module.Commands.Lint.Image == "" {
		return streamCommand(module.Commands.Lint.Run, module.Location)
	}

	env, err := environment.NewDockerEnvironment(
		module.Commands.Lint.Image,
		module.Location,
	)

	if err != nil {
		return err
	}

	ctx := context.Background()
	return env.Execute(ctx, module.Commands.Lint.Run)
}
