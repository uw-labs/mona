package command

import (
	"fmt"
	"path/filepath"

	"github.com/davidsbond/mona/internal/app"
	"github.com/davidsbond/mona/internal/config"
)

// AddApp creates a new "app.yaml" file in the specified directory
func AddApp(pj *config.Project, name, location string) error {
	apps, err := app.FindApps(pj.Location, pj.Mod)

	if err != nil {
		return err
	}

	for _, app := range apps {
		if name == app.Name {
			return fmt.Errorf("an app named '%s' already exists at %s", app.Name, app.Location)
		}
	}

	location = filepath.Join(pj.Location, location)

	if _, err := app.LoadApp(location, pj.Mod); err == app.ErrNoApp {
		return app.NewAppFile(name, location)
	} else if err != nil {
		return err
	}

	return fmt.Errorf("an app already exists in directory %s", location)
}
