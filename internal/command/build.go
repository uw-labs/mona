package command

import (
	"os"

	"github.com/davidsbond/mona/internal/files"
	"github.com/davidsbond/mona/internal/output"
)

// Build will execute the build commands for all modules where changes
// are detected.
func Build(pj *files.ProjectFile, parallelism int) error {
	return rangeChangedModules(pj, buildModule, rangeOptions{
		parallelism: parallelism,
		changeType:  changeTypeBuild,
	})
}

func buildModule(module *files.ModuleFile) error {
	if err := output.Writef(os.Stdout, "Building module %s", module.Name); err != nil {
		return err
	}

	if module.Commands.Build == "" {
		return nil
	}

	return streamCommand(module.Commands.Build, module.Location)
}
