package config_test

import (
	"testing"

	"github.com/davidsbond/mona/internal/config"
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
			if err := config.NewLockFile(".", tc.ProjectName); err != nil {
				assert.Fail(t, err.Error())
				return
			}

			defer deleteLockFile(t)
			assert.FileExists(t, "mona.lock")

			lock, err := config.LoadLockFile(".")

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
			if err := config.NewLockFile(".", tc.ProjectName); err != nil {
				assert.Fail(t, err.Error())
				return
			}

			defer deleteLockFile(t)
			assert.FileExists(t, "mona.lock")

			if err := config.UpdateLockFile(".", &config.LockFile{
				Name: tc.NewProjectName,
			}); err != nil {
				assert.Fail(t, err.Error())
				return
			}

			lock, err := config.LoadLockFile(".")

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
			if err := config.NewLockFile(".", tc.ProjectName); err != nil {
				assert.Fail(t, err.Error())
				return
			}

			defer deleteLockFile(t)
			assert.FileExists(t, "mona.lock")

			lock, err := config.LoadLockFile(".")

			if err != nil {
				assert.Fail(t, err.Error())
				return
			}

			assert.Equal(t, tc.ProjectName, lock.Name)
		})
	}
}

func TestAddApp(t *testing.T) {
	tt := []struct {
		Name        string
		ProjectName string
		AppName     string
		AppLocation string
		AppHash     string
	}{
		{
			Name:        "It should add an app to a lock file",
			ProjectName: "test",
			AppName:     "test",
			AppLocation: "test",
			AppHash:     "1234",
		},
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			if err := config.NewLockFile(".", tc.ProjectName); err != nil {
				assert.Fail(t, err.Error())
				return
			}

			defer deleteLockFile(t)
			assert.FileExists(t, "mona.lock")

			lock, err := config.LoadLockFile(".")

			if err != nil {
				assert.Fail(t, err.Error())
				return
			}

			if err := config.AddApp(lock, ".", tc.AppName); err != nil {
				assert.Fail(t, err.Error())
			}

			newLock, err := config.LoadLockFile(".")

			if err != nil {
				assert.Fail(t, err.Error())
				return
			}

			assert.Len(t, newLock.Apps, 1)
		})
	}
}
