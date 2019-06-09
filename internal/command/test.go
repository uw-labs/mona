package command

import (
	"github.com/davidsbond/mona/internal/files"
)

// Test attempts to run the test command for all modules where changes
// are detected.
func Test(pj *files.ProjectFile, parallelism int) error {
	return rangeChangedModules(pj, testModule, rangeOptions{
		updateHashes: true,
		changeType:   changeTypeTest,
		parallelism:  parallelism,
	})
}

func testModule(module *files.ModuleFile) error {
	if module.Commands.Test == "" {
		return nil
	}

	return streamCommand(module.Commands.Test, module.Location)
}
