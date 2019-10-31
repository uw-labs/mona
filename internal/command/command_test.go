package command_test

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/davidsbond/mona/internal/app"

	"github.com/davidsbond/mona/internal/config"

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

func deleteAppFile(t *testing.T, location string) {
	location = path.Join(location, "app.yml")
	if err := os.Remove(location); err != nil {
		assert.Fail(t, err.Error())
	}

}

func deleteAppFiles(t *testing.T, locations ...string) {
	for _, location := range locations {
		deleteAppFile(t, location)

		if err := os.RemoveAll(location); err != nil {
			assert.Fail(t, err.Error())
			return
		}
	}
}

func setupProject(t *testing.T) *config.Project {
	if err := command.Init(".", "test"); err != nil {
		assert.Fail(t, err.Error())
	}

	pj, err := config.LoadProject(".")

	if err != nil {
		assert.Fail(t, err.Error())
	}

	return pj
}

func setupApps(t *testing.T, pj *config.Project, locations ...string) {
	for _, location := range locations {
		if err := os.MkdirAll(location, 0777); err != nil {
			assert.Fail(t, err.Error())
		}

		if err := command.AddApp(pj, filepath.Base(location), location); err != nil {
			assert.Fail(t, err.Error())
		}
	}
}

func setupAppCode(t *testing.T, locations ...string) {
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

		mod, err := app.LoadApp(location)

		if err != nil {
			assert.FailNow(t, err.Error())
			return
		}

		mod.Commands.Build.Run = "go build"

		if err := app.SaveApp(location, mod); err != nil {
			assert.Fail(t, err.Error())
		}
	}
}
