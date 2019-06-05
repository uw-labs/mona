package files_test

import (
	"testing"

	"github.com/davidsbond/mona/internal/files"
	"github.com/stretchr/testify/assert"
)

func TestNewLockFile(t *testing.T) {
	tt := []struct {
		Name           string
		ProjectName    string
		ProjectVersion string
	}{
		{
			Name:           "It should create a lock file",
			ProjectName:    "test",
			ProjectVersion: "v1",
		},
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			if err := files.NewLockFile(tc.ProjectName, tc.ProjectVersion); err != nil {
				assert.Fail(t, err.Error())
				return
			}

			defer deleteLockFile(t)
			assert.FileExists(t, "mona.lock")

			lock, err := files.LoadLockFile(".")

			if err != nil {
				assert.Fail(t, err.Error())
				return
			}

			assert.Equal(t, tc.ProjectName, lock.Name)
			assert.Equal(t, tc.ProjectVersion, lock.Version)
		})
	}
}

func TestUpdateLockFile(t *testing.T) {
	tt := []struct {
		Name              string
		ProjectName       string
		ProjectVersion    string
		NewProjectName    string
		NewProjectVersion string
	}{
		{
			Name:              "It should update a lock file",
			ProjectName:       "test",
			ProjectVersion:    "v1",
			NewProjectName:    "test1",
			NewProjectVersion: "v2",
		},
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			if err := files.NewLockFile(tc.ProjectName, tc.ProjectVersion); err != nil {
				assert.Fail(t, err.Error())
				return
			}

			defer deleteLockFile(t)
			assert.FileExists(t, "mona.lock")

			if err := files.UpdateLockFile(&files.LockFile{
				Name:    tc.NewProjectName,
				Version: tc.NewProjectVersion,
			}); err != nil {
				assert.Fail(t, err.Error())
				return
			}

			lock, err := files.LoadLockFile(".")

			if err != nil {
				assert.Fail(t, err.Error())
				return
			}

			assert.Equal(t, tc.NewProjectName, lock.Name)
			assert.Equal(t, tc.NewProjectVersion, lock.Version)
		})
	}

}

func TestLoadLockFile(t *testing.T) {
	tt := []struct {
		Name           string
		ProjectName    string
		ProjectVersion string
	}{
		{
			Name:           "It should load a lock file",
			ProjectName:    "test",
			ProjectVersion: "v1",
		},
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			if err := files.NewLockFile(tc.ProjectName, tc.ProjectVersion); err != nil {
				assert.Fail(t, err.Error())
				return
			}

			defer deleteLockFile(t)
			assert.FileExists(t, "mona.lock")

			lock, err := files.LoadLockFile(".")

			if err != nil {
				assert.Fail(t, err.Error())
				return
			}

			assert.Equal(t, tc.ProjectName, lock.Name)
			assert.Equal(t, tc.ProjectVersion, lock.Version)
		})
	}
}

func TestLock_AddModule(t *testing.T) {
	tt := []struct {
		Name           string
		ProjectName    string
		ProjectVersion string
		ModuleName     string
		ModuleLocation string
		ModuleHash     string
	}{
		{
			Name:           "It should add a module to a lock file",
			ProjectName:    "test",
			ProjectVersion: "v1",
			ModuleName:     "test",
			ModuleLocation: "test",
			ModuleHash:     "1234",
		},
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			if err := files.NewLockFile(tc.ProjectName, tc.ProjectVersion); err != nil {
				assert.Fail(t, err.Error())
				return
			}

			defer deleteLockFile(t)
			assert.FileExists(t, "mona.lock")

			lock, err := files.LoadLockFile(".")

			if err != nil {
				assert.Fail(t, err.Error())
				return
			}

			if err := lock.AddModule(tc.ModuleName); err != nil {
				assert.Fail(t, err.Error())
			}

			newLock, err := files.LoadLockFile(".")

			if err != nil {
				assert.Fail(t, err.Error())
				return
			}

			assert.Len(t, newLock.Modules, 1)
		})
	}
}
