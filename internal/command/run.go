package command

import (
	"os"

	"github.com/apex/log"

	"github.com/davidsbond/mona/internal/app"
)

// Run iterates over all new/modified apps and executes their lint, test and build commands
// as necessary. Once complete, appropriate hashes are updated in the lock file.
func Run(cfg Config) error {
	if err := os.MkdirAll(cfg.Project.BinDir, os.ModePerm); err != nil {
		return err
	}

	lintChecked := make(map[string]bool)
	testChecked := make(map[string]bool)

	action := func(app *app.App, changes map[changeType]bool) error {
		if changes[changeTypeLint] {
			log.Infof("Linting app %s at %s", app.Name, app.Location)

			if err := app.Lint(lintChecked); err != nil {
				return err
			}
		}

		if changes[changeTypeTest] {
			log.Infof("Testing app %s at %s", app.Name, app.Location)

			if err := app.Test(testChecked); err != nil {
				return err
			}
		}

		if changes[changeTypeBuild] {
			log.Infof("Building app %s at %s", app.Name, app.Location)

			return app.Build("./" + cfg.Project.BinDir)
		}

		return nil
	}

	return rangeChangedApps(cfg, []changeType{changeTypeLint, changeTypeTest, changeTypeBuild}, action)
}

// Lint iterates over all new/modified apps and executes their lint command. Once complete,
// the lint hash is updated in the lock file.
func Lint(cfg Config) error {
	checked := make(map[string]bool)

	return rangeChangedApps(cfg, []changeType{changeTypeLint}, func(app *app.App, _ map[changeType]bool) error {
		log.Infof("Linting app %s at %s", app.Name, app.Location)

		return app.Lint(checked)
	})
}

// Test attempts to run the test command for all apps where changes
// are detected.
func Test(cfg Config) error {
	checked := make(map[string]bool)

	return rangeChangedApps(cfg, []changeType{changeTypeTest}, func(app *app.App, _ map[changeType]bool) error {
		log.Infof("Testing app %s at %s", app.Name, app.Location)

		return app.Test(checked)
	})
}

// Build will execute the build commands for all apps where changes
// are detected.
func Build(cfg Config) error {
	if err := os.MkdirAll(cfg.Project.BinDir, os.ModePerm); err != nil {
		return err
	}

	return rangeChangedApps(cfg, []changeType{changeTypeBuild}, func(app *app.App, _ map[changeType]bool) error {
		log.Infof("Building app %s at %s", app.Name, app.Location)

		return app.Build("./" + cfg.Project.BinDir)
	})
}
