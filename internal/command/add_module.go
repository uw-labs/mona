package command

import (
	"fmt"
	"path/filepath"

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
			return fmt.Errorf("a module named '%s' already exists at %s", module.Name, module.Location)
		}
	}

	location = filepath.Join(wd, location)

	if _, err := files.LoadModuleFile(location); err == files.ErrNoModule {
		return files.NewModuleFile(name, location)
	} else if err != nil {
		return err
	}

	return fmt.Errorf("a module already exists in directory %s", location)
}
