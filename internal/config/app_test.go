package config_test

import (
	"os"
	"testing"

	"github.com/davidsbond/mona/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestLoadAppFile(t *testing.T) {
	tt := []struct {
		Name     string
		AppName  string
		Expected *config.AppFile
	}{
		{
			Name:    "It should load an app file",
			AppName: "test",
			Expected: &config.AppFile{
				Name:     "test",
				Location: ".",
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			if err := config.NewAppFile(tc.AppName, "."); err != nil {
				assert.Fail(t, err.Error())
				return
			}

			mod, err := config.LoadAppFile(".")

			if err != nil {
				assert.Fail(t, err.Error())
				return
			}

			assert.EqualValues(t, tc.Expected, mod)

			if err := os.Remove("app.yml"); err != nil {
				assert.Fail(t, err.Error())
			}
		})
	}
}

func TestFindApps(t *testing.T) {
	tt := []struct {
		Name string
		Apps map[string]string
	}{
		{
			Name: "It should find all apps",
			Apps: map[string]string{
				"1": "./testdir/1",
				"2": "./testdir/2",
				"3": "./testdir/3",
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			for name, location := range tc.Apps {
				if err := os.MkdirAll(location, os.ModePerm); err != nil {
					assert.Fail(t, err.Error())
					return
				}

				if err := config.NewAppFile(name, location); err != nil {
					assert.Fail(t, err.Error())
					return
				}
			}

			apps, err := config.FindApps("./testdir")

			if err != nil {
				assert.Fail(t, err.Error())
				return
			}

			assert.Len(t, apps, len(tc.Apps))

			for name, location := range tc.Apps {
				mod, err := config.LoadAppFile(location)

				if err != nil {
					assert.Fail(t, err.Error())
					return
				}

				assert.Equal(t, name, mod.Name)
			}

			if err := os.RemoveAll("./testdir"); err != nil {
				assert.Fail(t, err.Error())
			}
		})
	}
}

func TestNewAppFile(t *testing.T) {
	tt := []struct {
		Name     string
		AppName  string
		Expected *config.AppFile
	}{
		{
			Name:    "It should create an app file",
			AppName: "test",
			Expected: &config.AppFile{
				Name:     "test",
				Location: ".",
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			if err := config.NewAppFile(tc.AppName, "."); err != nil {
				assert.Fail(t, err.Error())
				return
			}

			mod, err := config.LoadAppFile(".")

			if err != nil {
				assert.Fail(t, err.Error())
				return
			}

			assert.EqualValues(t, tc.Expected, mod)

			if err := os.Remove("app.yml"); err != nil {
				assert.Fail(t, err.Error())
			}
		})
	}
}

func TestUpdateAppFile(t *testing.T) {
	tt := []struct {
		Name       string
		AppName    string
		NewAppData *config.AppFile
	}{
		{
			Name:    "It should update an app file",
			AppName: "test",
			NewAppData: &config.AppFile{
				Name:     "test1",
				Location: ".",
				Exclude:  []string{"test"},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			if err := config.NewAppFile(tc.AppName, "."); err != nil {
				assert.Fail(t, err.Error())
				return
			}

			if err := config.UpdateAppFile(".", tc.NewAppData); err != nil {
				assert.Fail(t, err.Error())
				return
			}

			mod, err := config.LoadAppFile(".")

			if err != nil {
				assert.Fail(t, err.Error())
				return
			}

			assert.EqualValues(t, tc.NewAppData, mod)

			if err := os.Remove("app.yml"); err != nil {
				assert.Fail(t, err.Error())
			}
		})
	}
}
