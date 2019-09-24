package command

import (
	"context"

	"github.com/apex/log"
	"github.com/davidsbond/mona/internal/config"
	"github.com/davidsbond/mona/internal/deps"
	"github.com/davidsbond/mona/internal/environment"
)

// Build will execute the build commands for all apps where changes
// are detected.
func Build(mod deps.Module, pj *config.ProjectFile) error {
	return rangeChangedApps(mod, pj, buildApp, changeTypeBuild)
}

func buildApp(app *config.AppFile) error {
	log.Infof("Building app %s at %s", app.Name, app.Location)

	if app.Commands.Build.Run == "" {
		return nil
	}

	// Run command locally if no image is set
	if app.Commands.Build.Image == "" {
		return streamCommand(app.Commands.Build.Run, app.Location)
	}

	env, err := environment.NewDockerEnvironment(
		app.Commands.Build.Image,
		app.Location,
	)

	if err != nil {
		return err
	}

	ctx := context.Background()
	return env.Execute(ctx, app.Commands.Build.Run)
}
