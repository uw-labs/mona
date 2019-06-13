package command

import (
	"github.com/apex/log"

	"github.com/davidsbond/mona/internal/files"
)

// Test attempts to run the test command for all modules where changes
// are detected.
func Test(pj *files.ProjectFile) error {
	return rangeChangedModules(pj, testModule, changeTypeTest)
}

func testModule(module *files.ModuleFile) error {
	log.Infof("Testing module %s at %s", module.Name, module.Location)

	if module.Commands.Test == "" {
		return nil
	}

	return streamCommand(module.Commands.Test, module.Location)
}
