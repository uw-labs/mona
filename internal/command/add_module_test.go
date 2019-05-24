package command_test

import (
	"testing"

	"github.com/davidsbond/mona/internal/command"
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
			setupProject(t)
			defer deleteProjectFiles(t)
			defer deleteModuleFile(t, tc.ModuleLocation)

			if err := command.AddModule(tc.ModuleName, tc.ModuleLocation); err != nil {
				assert.FailNow(t, err.Error())
			}

			assert.FileExists(t, "module.yml")
		})
	}
}
