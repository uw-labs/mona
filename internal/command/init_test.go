package command_test

import (
	"testing"

	"github.com/davidsbond/mona/internal/command"
	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	tt := []struct {
		Name           string
		ProjectName    string
		ProjectVersion string
	}{
		{
			Name:           "It should create a mona.yaml and mona.lock file",
			ProjectName:    "test",
			ProjectVersion: "test",
		},
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			defer deleteProjectFiles(t)

			if err := command.Init(tc.ProjectName, tc.ProjectVersion); err != nil {
				assert.Fail(t, err.Error())
			}

			assert.FileExists(t, "mona.yml")
			assert.FileExists(t, "mona.lock")
		})
	}
}
