package command

import (
	"os/exec"
	"strings"

	"github.com/davidsbond/mona/internal/files"
)

// Test attempts to run the test command for all modules where changes
// are detected.
func Test(pj *files.ProjectFile) error {
	return rangeChangedModules(pj, changeTypeTest, testModule, true)
}

func testModule(module *files.ModuleFile) error {
	if module.Commands.Test == "" {
		return nil
	}

	parts := strings.Split(module.Commands.Test, " ")
	cmd := exec.Command(parts[0], parts[1:]...)
	cmd.Dir = module.Location

	return streamCommand(module.Commands.Test, module.Location)
}
