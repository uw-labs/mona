package files_test

import (
	"testing"

	"github.com/davidsbond/mona/internal/files"
	"github.com/stretchr/testify/assert"
)

func TestNewLockFile(t *testing.T) {
	tt := []struct {
		Name        string
		ProjectName string
	}{
		{
			Name:        "It should create a lock file",
			ProjectName: "test",
		},
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			if err := files.NewLockFile(".", tc.ProjectName); err != nil {
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
		})
	}
}

func TestUpdateLockFile(t *testing.T) {
	tt := []struct {
		Name           string
		ProjectName    string
		NewProjectName string
	}{
		{
			Name:           "It should update a lock file",
			ProjectName:    "test",
			NewProjectName: "test1",
		},
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			if err := files.NewLockFile(".", tc.ProjectName); err != nil {
				assert.Fail(t, err.Error())
				return
			}

			defer deleteLockFile(t)
			assert.FileExists(t, "mona.lock")

			if err := files.UpdateLockFile(".", &files.LockFile{
				Name: tc.NewProjectName,
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
		})
	}

}

func TestLoadLockFile(t *testing.T) {
	tt := []struct {
		Name        string
		ProjectName string
	}{
		{
			Name:        "It should load a lock file",
			ProjectName: "test",
		},
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			if err := files.NewLockFile(".", tc.ProjectName); err != nil {
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
		})
	}
}

func TestAddModule(t *testing.T) {
	tt := []struct {
		Name           string
		ProjectName    string
		ModuleName     string
		ModuleLocation string
		ModuleHash     string
	}{
		{
			Name:           "It should add a module to a lock file",
			ProjectName:    "test",
			ModuleName:     "test",
			ModuleLocation: "test",
			ModuleHash:     "1234",
		},
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			if err := files.NewLockFile(".", tc.ProjectName); err != nil {
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

			if err := files.AddModule(lock, ".", tc.ModuleName); err != nil {
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
