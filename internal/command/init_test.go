package command_test

import (
	"testing"

	"github.com/davidsbond/mona/internal/command"
	"github.com/davidsbond/mona/internal/files"
	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	tt := []struct {
		Name        string
		ProjectName string
	}{
		{
			Name:        "It should create a mona.yaml and mona.lock file",
			ProjectName: "test",
		},
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			defer deleteProjectFiles(t)

			if err := command.Init(tc.ProjectName); err != nil {
				assert.Fail(t, err.Error())
				return
			}

			assert.FileExists(t, "mona.yml")
			assert.FileExists(t, "mona.lock")

			proj, err := files.LoadProjectFile(".")

			if err != nil {
				assert.Fail(t, err.Error())
				return
			}

			assert.Equal(t, tc.ProjectName, proj.Name)

			lock, err := files.LoadLockFile(".")

			if err != nil {
				assert.Fail(t, err.Error())
				return
			}

			assert.Equal(t, tc.ProjectName, lock.Name)
		})
	}
}
