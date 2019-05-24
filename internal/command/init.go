package command

import "github.com/davidsbond/mona/internal/files"

// Init creates a new project and lock file using the provided name and
// version
func Init(name, version string) error {
	if err := files.NewProjectFile(name, version); err != nil {
		return err
	}

	return files.NewLockFile(name, version)
}
