package command_test

import (
	"testing"

	"github.com/davidsbond/mona/internal/command"
	"github.com/stretchr/testify/assert"
)

func TestDiff(t *testing.T) {
	tt := []struct {
		Name          string
		ModuleDirs    []string
		ExpectedDiffs []string
	}{
		{
			Name:          "It should detect changes in modules",
			ModuleDirs:    []string{"test/a", "test/b"},
			ExpectedDiffs: []string{"a"},
		},
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			setupProject(t)
			setupModules(t, tc.ModuleDirs...)

			defer deleteProjectFiles(t)
			defer deleteModuleFiles(t, tc.ModuleDirs...)

			names, err := command.Diff()

			if err != nil {
				assert.Fail(t, err.Error())
				return
			}

			for _, exp := range tc.ExpectedDiffs {
				assert.Contains(t, names, exp)
			}
		})
	}
}
