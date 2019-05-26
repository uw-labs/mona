package command

import (
	"github.com/davidsbond/mona/internal/files"
)

// Build will execute the build commands for all modules where changes
// are detected.
func Build() error {
	return rangeChangedModules(changeTypeBuild, true, buildModule)
}

func buildModule(module *files.ModuleFile) error {
	if module.Commands.Build == "" {
		return nil
	}

	if err := streamCommand(module.Commands.Build, module.Location); err != nil {
		return err
	}

	if len(module.Artefacts) > 0 {
		project, err := files.LoadProjectFile()

		if err != nil {
			return err
		}

		if err := module.CollectArtefacts(project.Artefacts); err != nil {
			return err
		}
	}

	return nil
}
