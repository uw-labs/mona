package command_test

import (
	"testing"

	"github.com/davidsbond/mona/internal/command"
	"github.com/stretchr/testify/assert"
)

func TestBuild(t *testing.T) {
	tt := []struct {
		Name              string
		ModuleDirs        []string
		ExpectedArtefacts []string
	}{
		{
			Name:              "It should build all new modules",
			ModuleDirs:        []string{"test/a", "test/b"},
			ExpectedArtefacts: []string{"test/a/a", "test/b/b"},
		},
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			pj := setupProject(t)
			setupModules(t, pj, tc.ModuleDirs...)
			setupModuleCode(t, tc.ModuleDirs...)

			defer deleteModuleFiles(t, tc.ModuleDirs...)
			defer deleteProjectFiles(t)

			if err := command.Build(pj, 1); err != nil {
				assert.Fail(t, err.Error())
				return
			}

			for _, exp := range tc.ExpectedArtefacts {
				assert.FileExists(t, exp)
			}
		})
	}
}
