package command

import (
	"context"

	"github.com/apex/log"
	"github.com/davidsbond/mona/internal/config"
	"github.com/davidsbond/mona/internal/deps"
	"github.com/davidsbond/mona/internal/environment"
)

// Test attempts to run the test command for all apps where changes
// are detected.
func Test(mod deps.Module, pj *config.ProjectFile) error {
	return rangeChangedApps(mod, pj, testApp, changeTypeTest)
}

func testApp(app *config.AppFile) error {
	log.Infof("Testing app %s at %s", app.Name, app.Location)

	if app.Commands.Test.Run == "" {
		return nil
	}

	// Run command locally if no image is set
	if app.Commands.Test.Image == "" {
		return streamCommand(app.Commands.Test.Run, app.Location)
	}

	env, err := environment.NewDockerEnvironment(
		app.Commands.Test.Image,
		app.Location,
	)

	if err != nil {
		return err
	}

	ctx := context.Background()
	return env.Execute(ctx, app.Commands.Test.Run)
}
