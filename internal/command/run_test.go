package command_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/uw-labs/mona/internal/command"
)

func TestBuild(t *testing.T) {
	tt := []struct {
		Name              string
		AppDirs           []string
		ExpectedArtifacts []string
	}{
		{
			Name:              "It should build all new apps",
			AppDirs:           []string{"test/a", "test/b"},
			ExpectedArtifacts: []string{"test/a/a", "test/b/b"},
		},
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			pj := setupProject(t)
			setupApps(t, pj, tc.AppDirs...)
			setupAppCode(t, tc.AppDirs...)

			defer deleteAppFiles(t, tc.AppDirs...)
			defer deleteProjectFiles(t)

			if err := command.Build(golang.Module{}, pj); err != nil {
				assert.Fail(t, err.Error())
				return
			}

			for _, exp := range tc.ExpectedArtifacts {
				assert.FileExists(t, exp)
			}
		})
	}
}
