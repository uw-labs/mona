package command_test

import (
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/davidsbond/mona/internal/command"
	"github.com/stretchr/testify/assert"
)

func deleteProjectFiles(t *testing.T) {
	if err := os.Remove("mona.yml"); err != nil {
		assert.FailNow(t, err.Error())
	}

	if err := os.Remove("mona.lock"); err != nil {
		assert.FailNow(t, err.Error())
	}
}

func deleteModuleFile(t *testing.T, location string) {
	location = path.Join(location, "module.yml")
	if err := os.Remove(location); err != nil {
		assert.FailNow(t, err.Error())
	}

}

func deleteModuleFiles(t *testing.T, locations ...string) {
	for _, location := range locations {
		deleteModuleFile(t, location)

		if err := os.RemoveAll(location); err != nil {
			assert.FailNow(t, err.Error())
			return
		}
	}
}

func setupProject(t *testing.T) {
	if err := command.Init("test", "v1"); err != nil {
		assert.FailNow(t, err.Error())
	}
}

func setupModules(t *testing.T, locations ...string) {
	for _, location := range locations {
		if err := os.MkdirAll(location, 0777); err != nil {
			assert.FailNow(t, err.Error())
		}

		if err := command.AddModule(filepath.Base(location), location); err != nil {
			assert.FailNow(t, err.Error())
		}
	}
}
