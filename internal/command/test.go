package command

import (
	"github.com/apex/log"

	"github.com/davidsbond/mona/internal/app"
	"github.com/davidsbond/mona/internal/config"
)

// Test attempts to run the test command for all apps where changes
// are detected.
func Test(pj *config.ProjectFile) error {
	checked := make(map[string]bool)

	return rangeChangedApps(pj, changeTypeTest, func(app *app.App) error {
		log.Infof("Testing app %s at %s", app.Name, app.Location)

		return app.Test(checked)
	})
}
