package command_test

import (
	"testing"

	"github.com/davidsbond/mona/internal/command"
	"github.com/davidsbond/mona/internal/deps"
	"github.com/stretchr/testify/assert"
)

func TestDiff(t *testing.T) {
	tt := []struct {
		Name          string
		AppDirs       []string
		ExpectedDiffs []string
	}{
		{
			Name:          "It should detect changes in apps",
			AppDirs:       []string{"test/a", "test/b"},
			ExpectedDiffs: []string{"a"},
		},
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			pj := setupProject(t)
			setupApps(t, pj, tc.AppDirs...)

			defer deleteProjectFiles(t)
			defer deleteAppFiles(t, tc.AppDirs...)

			build, test, lint, err := command.Diff(deps.Module{}, pj)

			if err != nil {
				assert.Fail(t, err.Error())
				return
			}

			for _, exp := range tc.ExpectedDiffs {
				assert.Contains(t, build, exp)
				assert.Contains(t, test, exp)
				assert.Contains(t, lint, exp)
			}
		})
	}
}
