package walk_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/davidsbond/mona/pkg/walk"
	"github.com/stretchr/testify/assert"
)

func TestFast(t *testing.T) {
	t.Parallel()

	tt := []struct {
		Name          string
		Directory     string
		ExpectedFiles []string
		Parallelism   int
	}{
		{
			Name:      "It should walk over all files in the directory",
			Directory: "test",
			ExpectedFiles: []string{
				filepath.Join("test", "1.txt"),
				filepath.Join("test", "2.txt"),
				filepath.Join("test", "3.txt"),
			},
			Parallelism: 1,
		},
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			setupTestDirectory(t, tc.Directory, tc.ExpectedFiles)

			err := walk.Fast(tc.Directory, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					assert.Fail(t, err.Error())
					return nil
				}

				if info.IsDir() {
					return nil
				}

				assert.Contains(t, tc.ExpectedFiles, path)
				return nil
			}, tc.Parallelism)

			if err != nil {
				assert.Fail(t, err.Error())
				return
			}

			if err := os.RemoveAll(tc.Directory); err != nil {
				assert.Fail(t, err.Error())
			}
		})
	}
}

func setupTestDirectory(t *testing.T, dir string, files []string) {
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		assert.FailNow(t, err.Error())
		return
	}

	for _, file := range files {
		path, _ := filepath.Split(file)

		if err := os.MkdirAll(path, os.ModePerm); err != nil {
			assert.FailNow(t, err.Error())
			return
		}

		file, err := os.Create(file)

		if err != nil {
			assert.FailNow(t, err.Error())
			return
		}

		if err := file.Close(); err != nil {
			assert.FailNow(t, err.Error())
			return
		}
	}
}
