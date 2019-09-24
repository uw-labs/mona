package command_test

import (
	"testing"

	"github.com/davidsbond/mona/internal/command"
	"github.com/davidsbond/mona/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestAddApp(t *testing.T) {
	tt := []struct {
		Name        string
		AppName     string
		AppLocation string
	}{
		{
			Name:        "It should create an app",
			AppName:     "test",
			AppLocation: ".",
		},
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			pj := setupProject(t)
			defer deleteProjectFiles(t)
			defer deleteAppFile(t, tc.AppLocation)

			if err := command.AddApp(pj, tc.AppName, tc.AppLocation); err != nil {
				assert.Fail(t, err.Error())
				return
			}

			assert.FileExists(t, "app.yml")
			mod, err := config.LoadAppFile(tc.AppLocation)

			if err != nil {
				assert.Fail(t, err.Error())
				return
			}

			assert.Equal(t, tc.AppName, mod.Name)
		})
	}
}
