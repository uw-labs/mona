package files_test

import (
	"testing"

	"github.com/davidsbond/mona/internal/files"
	"github.com/stretchr/testify/assert"
)

func TestNewProjectFile(t *testing.T) {
	tt := []struct {
		Name           string
		ProjectName    string
		ProjectVersion string
	}{
		{
			Name:           "It should create a project file",
			ProjectName:    "test",
			ProjectVersion: "v1",
		},
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			if err := files.NewProjectFile(tc.ProjectName, tc.ProjectVersion); err != nil {
				assert.Fail(t, err.Error())
				return
			}

			defer deleteProjectFile(t)

			assert.FileExists(t, "mona.yml")
			proj, err := files.LoadProjectFile()

			if err != nil {
				assert.Fail(t, err.Error())
				return
			}

			assert.Equal(t, tc.ProjectName, proj.Name)
			assert.Equal(t, tc.ProjectVersion, proj.Version)
		})
	}
}

func TestLoadProjectFile(t *testing.T) {
	tt := []struct {
		Name           string
		ProjectName    string
		ProjectVersion string
	}{
		{
			Name:           "It should create a project file",
			ProjectName:    "test",
			ProjectVersion: "v1",
		},
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			if err := files.NewProjectFile(tc.ProjectName, tc.ProjectVersion); err != nil {
				assert.Fail(t, err.Error())
				return
			}

			defer deleteProjectFile(t)
			assert.FileExists(t, "mona.yml")

			proj, err := files.LoadProjectFile()

			if err != nil {
				assert.Fail(t, err.Error())
				return
			}

			assert.Equal(t, tc.ProjectName, proj.Name)
			assert.Equal(t, tc.ProjectVersion, proj.Version)
		})
	}
}

func TestProject_AddModule(t *testing.T) {
	tt := []struct {
		Name string
	}{}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {

		})
	}
}
