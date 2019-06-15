package command

import (
	"context"

	"github.com/apex/log"
	"github.com/davidsbond/mona/internal/environment"
	"github.com/davidsbond/mona/internal/files"
)

// Test attempts to run the test command for all modules where changes
// are detected.
func Test(pj *files.ProjectFile) error {
	return rangeChangedModules(pj, testModule, changeTypeTest)
}

func testModule(module *files.ModuleFile) error {
	log.Infof("Testing module %s at %s", module.Name, module.Location)

	if module.Commands.Test.Run == "" {
		return nil
	}

	// Run command locally if no image is set
	if module.Commands.Test.Image == "" {
		return streamCommand(module.Commands.Test.Run, module.Location)
	}

	env, err := environment.NewDockerEnvironment(
		module.Commands.Test.Image,
		module.Location,
	)

	if err != nil {
		return err
	}

	ctx := context.Background()
	return env.Execute(ctx, module.Commands.Test.Run)
}
