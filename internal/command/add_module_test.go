package command_test

import (
	"testing"

	"github.com/davidsbond/mona/internal/command"
	"github.com/davidsbond/mona/internal/files"
	"github.com/stretchr/testify/assert"
)

func TestAddModule(t *testing.T) {
	tt := []struct {
		Name           string
		ModuleName     string
		ModuleLocation string
	}{
		{
			Name:           "It should create a module",
			ModuleName:     "test",
			ModuleLocation: ".",
		},
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			pj := setupProject(t)
			defer deleteProjectFiles(t)
			defer deleteModuleFile(t, tc.ModuleLocation)

			if err := command.AddModule(pj, tc.ModuleName, tc.ModuleLocation); err != nil {
				assert.Fail(t, err.Error())
				return
			}

			assert.FileExists(t, "module.yml")
			mod, err := files.LoadModuleFile(tc.ModuleLocation)

			if err != nil {
				assert.Fail(t, err.Error())
				return
			}

			assert.Equal(t, tc.ModuleName, mod.Name)
		})
	}
}
