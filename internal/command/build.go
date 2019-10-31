package command

import (
	"os"

	"github.com/apex/log"

	"github.com/davidsbond/mona/internal/app"
	"github.com/davidsbond/mona/internal/config"
)

// Build will execute the build commands for all apps where changes
// are detected.
func Build(pj *config.Project) error {
	return rangeChangedApps(pj, changeTypeBuild, func(app *app.App) error {
		log.Infof("Building app %s at %s", app.Name, app.Location)
		if err := os.MkdirAll(pj.BinDir, os.ModePerm); err != nil {
			return err
		}

		return app.Build("./" + pj.BinDir)
	})
}
