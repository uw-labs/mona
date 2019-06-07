package command_test

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/davidsbond/mona/internal/files"

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

func deleteModuleFiles(t *testing.T, locations ...string) {
	for _, location := range locations {
		deleteModuleFile(t, location)

		if err := os.RemoveAll(location); err != nil {
			assert.Fail(t, err.Error())
			return
		}
	}
}

func setupProject(t *testing.T) *files.ProjectFile {
	if err := command.Init("test"); err != nil {
		assert.Fail(t, err.Error())
	}

	pj, err := files.LoadProjectFile(".")

	if err != nil {
		assert.Fail(t, err.Error())
	}

	return pj
}

func setupModules(t *testing.T, pj *files.ProjectFile, locations ...string) {
	for _, location := range locations {
		if err := os.MkdirAll(location, 0777); err != nil {
			assert.Fail(t, err.Error())
		}

		if err := command.AddModule(pj, filepath.Base(location), location); err != nil {
			assert.Fail(t, err.Error())
		}
	}
}

func setupModuleCode(t *testing.T, locations ...string) {
	for _, location := range locations {
		src := filepath.Join(location, "main.go")
		data := `
		package main

		func main() {

		}
		`

		if err := ioutil.WriteFile(src, []byte(data), 0777); err != nil {
			assert.Fail(t, err.Error())
		}

		mod, err := files.LoadModuleFile(location)

		if err != nil {
			assert.FailNow(t, err.Error())
			return
		}

		mod.Commands.Build = "go build"

		if err := files.UpdateModuleFile(location, mod); err != nil {
			assert.Fail(t, err.Error())
		}
	}
}
