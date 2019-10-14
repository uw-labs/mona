package app_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/davidsbond/mona/internal/app"
)

func TestLoadAppFile(t *testing.T) {
	tt := []struct {
		Name     string
		AppName  string
		Expected *app.App
	}{
		{
			Name:    "It should load an app file",
			AppName: "test",
			Expected: &app.App{
				Name:     "test",
				Location: ".",
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			if err := app.NewAppFile(tc.AppName, "."); err != nil {
				assert.Fail(t, err.Error())
				return
			}

			mod, err := app.LoadApp(".")

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

				if err := app.NewAppFile(name, location); err != nil {
					assert.Fail(t, err.Error())
					return
				}
			}

			apps, err := app.FindApps("./testdir")

			if err != nil {
				assert.Fail(t, err.Error())
				return
			}

			assert.Len(t, apps, len(tc.Apps))

			for name, location := range tc.Apps {
				mod, err := app.LoadApp(location)

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
		Expected *app.App
	}{
		{
			Name:    "It should create an app file",
			AppName: "test",
			Expected: &app.App{
				Name:     "test",
				Location: ".",
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			if err := app.NewAppFile(tc.AppName, "."); err != nil {
				assert.Fail(t, err.Error())
				return
			}

			mod, err := app.LoadApp(".")

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
		NewAppData *app.App
	}{
		{
			Name:    "It should update an app file",
			AppName: "test",
			NewAppData: &app.App{
				Name:     "test1",
				Location: ".",
				Exclude:  []string{"test"},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			if err := app.NewAppFile(tc.AppName, "."); err != nil {
				assert.Fail(t, err.Error())
				return
			}

			if err := app.SaveApp(".", tc.NewAppData); err != nil {
				assert.Fail(t, err.Error())
				return
			}

			mod, err := app.LoadApp(".")

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
