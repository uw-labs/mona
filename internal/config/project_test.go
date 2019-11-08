package config_test

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/uw-labs/mona/internal/config"
)

var goModData = []byte(`module github.com/some/project

go 1.12

require (
	github.com/urfave/cli v1.20.0
)

`)

func deleteProjectFile(t *testing.T) {
	require.NoError(t, os.Remove("mona.yml"))
}

func createGoMod(t *testing.T) {
	require.NoError(t, ioutil.WriteFile("go.mod", goModData, os.ModePerm))
}

func deleteGoMod(t *testing.T) {
	require.NoError(t, os.Remove("go.mod"))
}

func TestProjectCreateAndLoad(t *testing.T) {
	createGoMod(t)
	defer deleteGoMod(t)

	require.NoError(t, config.NewProject(".", "test"))
	defer deleteProjectFile(t)

	require.FileExists(t, "mona.yml")

	project, err := config.LoadProject(".")
	require.NoError(t, err)

	expected := &config.Project{
		Name:     "test",
		Location: ".",
		BinDir:   "bin",
		Mod: golang.Module{
			Name: "github.com/some/project",
			Deps: map[string]string{
				"github.com/urfave/cli": "v1.20.0",
			},
		},
	}

	require.Equal(t, expected, project)
}

func TestNewProject_NoGoModFile(t *testing.T) {
	require.Equal(t, config.ErrNoGoModFile, config.NewProject(".", "test"))
}
