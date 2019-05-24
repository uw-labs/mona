package command_test

import (
	"os"
	"path"
	"testing"

	"github.com/davidsbond/mona/internal/command"
	"github.com/stretchr/testify/assert"
)

func deleteProjectFiles(t *testing.T) {
	if err := os.Remove("mona.yml"); err != nil {
		assert.Fail(t, err.Error())
	}

	if err := os.Remove("mona.lock"); err != nil {
		assert.Fail(t, err.Error())
	}
}

func deleteModuleFile(t *testing.T, location string) {
	location = path.Join(location, "module.yml")
	if err := os.Remove(location); err != nil {
		assert.Fail(t, err.Error())
	}

}

func setupProject(t *testing.T) {
	if err := command.Init("test", "v1"); err != nil {
		assert.Fail(t, err.Error())
	}
}
