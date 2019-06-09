package files_test

import (
	"os"
	"testing"

	"github.com/davidsbond/mona/internal/files"
	"github.com/stretchr/testify/assert"
)

func TestLoadModuleFile(t *testing.T) {
	tt := []struct {
		Name       string
		ModuleName string
		Expected   *files.ModuleFile
	}{
		{
			Name:       "It should load a module file",
			ModuleName: "test",
			Expected: &files.ModuleFile{
				Name:     "test",
				Location: ".",
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			if err := files.NewModuleFile(tc.ModuleName, "."); err != nil {
				assert.Fail(t, err.Error())
				return
			}

			mod, err := files.LoadModuleFile(".")

			if err != nil {
				assert.Fail(t, err.Error())
				return
			}

			assert.EqualValues(t, tc.Expected, mod)

			if err := os.Remove("module.yml"); err != nil {
				assert.Fail(t, err.Error())
			}
		})
	}
}

func TestFindModules(t *testing.T) {

}

func TestNewModuleFile(t *testing.T) {
	tt := []struct {
		Name       string
		ModuleName string
		Expected   *files.ModuleFile
	}{
		{
			Name:       "It should create a module file",
			ModuleName: "test",
			Expected: &files.ModuleFile{
				Name:     "test",
				Location: ".",
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			if err := files.NewModuleFile(tc.ModuleName, "."); err != nil {
				assert.Fail(t, err.Error())
				return
			}

			mod, err := files.LoadModuleFile(".")

			if err != nil {
				assert.Fail(t, err.Error())
				return
			}

			assert.EqualValues(t, tc.Expected, mod)

			if err := os.Remove("module.yml"); err != nil {
				assert.Fail(t, err.Error())
			}
		})
	}
}

func TestUpdateModuleFile(t *testing.T) {
	tt := []struct {
		Name          string
		ModuleName    string
		NewModuleData *files.ModuleFile
	}{
		{
			Name:       "It should update a module file",
			ModuleName: "test",
			NewModuleData: &files.ModuleFile{
				Name:     "test1",
				Location: ".",
				Exclude:  []string{"test"},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			if err := files.NewModuleFile(tc.ModuleName, "."); err != nil {
				assert.Fail(t, err.Error())
				return
			}

			if err := files.UpdateModuleFile(".", tc.NewModuleData); err != nil {
				assert.Fail(t, err.Error())
				return
			}

			mod, err := files.LoadModuleFile(".")

			if err != nil {
				assert.Fail(t, err.Error())
				return
			}

			assert.EqualValues(t, tc.NewModuleData, mod)

			if err := os.Remove("module.yml"); err != nil {
				assert.Fail(t, err.Error())
			}
		})
	}
}
