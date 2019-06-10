package command

import (
	"os"

	"github.com/davidsbond/mona/internal/files"
	"github.com/davidsbond/mona/internal/output"
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
	if err := output.Writef(os.Stdout, "Testing module %s", module.Name); err != nil {
		return err
	}

	if module.Commands.Test == "" {
		return nil
	}

	return streamCommand(module.Commands.Test, module.Location)
}
