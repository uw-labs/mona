package command

import (
	"context"

	"github.com/apex/log"
	"github.com/davidsbond/mona/internal/config"
	"github.com/davidsbond/mona/internal/deps"
	"github.com/davidsbond/mona/internal/environment"
)

// Lint iterates over all new/modified apps and executes their lint command. Once complete,
// the lint hash is updated in the lock file.
func Lint(mod deps.Module, pj *config.ProjectFile) error {
	return rangeChangedApps(mod, pj, lintApp, changeTypeLint)
}

func lintApp(app *config.AppFile) error {
	log.Infof("Linting app %s at %s", app.Name, app.Location)

	if app.Commands.Lint.Run == "" {
		return nil
	}

	// Run command locally if no image is specified
	if app.Commands.Lint.Image == "" {
		return streamCommand(app.Commands.Lint.Run, app.Location)
	}

	env, err := environment.NewDockerEnvironment(
		app.Commands.Lint.Image,
		app.Location,
	)

	if err != nil {
		return err
	}

	ctx := context.Background()
	return env.Execute(ctx, app.Commands.Lint.Run)
}
