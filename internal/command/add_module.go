package command

import (
	"fmt"

	"github.com/davidsbond/mona/internal/files"
)

// AddModule creates a new "module.yaml" file in the specified directory
func AddModule(wd, name, location string) error {
	modules, err := files.FindModules(wd)

	if err != nil {
		return err
	}

	for _, module := range modules {
		if name == module.Name {
			return fmt.Errorf("A module named '%s' already exists at %s", module.Name, module.Location)
		}
	}

	if _, err := files.LoadModuleFile(location); err == files.ErrNoModule {
		return files.NewModuleFile(name, location)
	} else if err != nil {
		return err
	}

	return fmt.Errorf("A module already exists in directory %s", location)
}
