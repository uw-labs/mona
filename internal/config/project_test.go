package config_test

import (
	"testing"

	"github.com/davidsbond/mona/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestNewProjectFile(t *testing.T) {
	tt := []struct {
		Name        string
		ProjectName string
	}{
		{
			Name:        "It should create a project file",
			ProjectName: "test",
		},
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			if err := config.NewProjectFile(".", tc.ProjectName); err != nil {
				assert.Fail(t, err.Error())
				return
			}

			defer deleteProjectFile(t)

			assert.FileExists(t, "mona.yml")
			proj, err := config.LoadProjectFile(".")

			if err != nil {
				assert.Fail(t, err.Error())
				return
			}

			assert.Equal(t, tc.ProjectName, proj.Name)
		})
	}
}

func TestLoadProjectFile(t *testing.T) {
	tt := []struct {
		Name        string
		ProjectName string
	}{
		{
			Name:        "It should create a project file",
			ProjectName: "test",
		},
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			if err := config.NewProjectFile(".", tc.ProjectName); err != nil {
				assert.Fail(t, err.Error())
				return
			}

			defer deleteProjectFile(t)
			assert.FileExists(t, "mona.yml")

			proj, err := config.LoadProjectFile(".")

			if err != nil {
				assert.Fail(t, err.Error())
				return
			}

			assert.Equal(t, tc.ProjectName, proj.Name)
		})
	}
}
