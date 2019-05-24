package command

import "github.com/davidsbond/mona/internal/files"

// AddModule creates a new "module.yaml" file in the specified directory and adds it to the project
// and lock files.
func AddModule(name, location string) error {
	if err := files.NewModuleFile(name, location); err != nil {
		return err
	}

	pj, err := files.LoadProjectFile()

	if err != nil {
		return err
	}

	if err := pj.AddModule(name, location); err != nil {
		return err
	}

	lock, err := files.LoadLockFile()

	if err != nil {
		return err
	}

	return lock.AddModule(name, location, "")
}
