package command

import "github.com/davidsbond/mona/internal/files"

// Init creates a new project and lock file in the provided working directory
// with the given name.
func Init(wd, name string) error {
	if err := files.NewProjectFile(wd, name); err != nil {
		return err
	}

	return files.NewLockFile(wd, name)
}
