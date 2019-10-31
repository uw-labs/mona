package command

import (
	"github.com/apex/log"

	"github.com/davidsbond/mona/internal/app"
	"github.com/davidsbond/mona/internal/config"
)

// Lint iterates over all new/modified apps and executes their lint command. Once complete,
// the lint hash is updated in the lock file.
func Lint(pj *config.Project) error {
	checked := make(map[string]bool)

	return rangeChangedApps(pj, changeTypeLint, func(app *app.App) error {
		log.Infof("Linting app %s at %s", app.Name, app.Location)

		return app.Lint(checked)
	})
}
