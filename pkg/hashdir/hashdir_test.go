package hashdir_test

import (
	"encoding/base64"
	"os"
	"testing"

	"github.com/davidsbond/mona/pkg/hashdir"
	"github.com/stretchr/testify/assert"
)

func TestGenerate(t *testing.T) {
	tt := []struct {
		Name        string
		Directory   string
		Parallelism int
		Expected    string
	}{
		{
			Name:        "It should generate a base64 hash of a directory",
			Directory:   "./testdir",
			Parallelism: 1,
			Expected:    "XrY7u+Ae7tCTyyK7j1rNww==",
		},
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			setup(t)
			defer teardown(t)

			hash, err := hashdir.Generate(tc.Directory, tc.Parallelism)

			if err != nil {
				assert.Fail(t, err.Error())
				return
			}

			assert.Equal(t, tc.Expected, hash)

			if _, err := base64.StdEncoding.Decode(make([]byte, 1024), []byte(hash)); err != nil {
				assert.Fail(t, err.Error())
			}
		})
	}
}

func setup(t *testing.T) {
	if err := os.Mkdir("./testdir", os.ModePerm); err != nil {
		assert.FailNow(t, err.Error())
		return
	}

	file, err := os.Create("./testdir/data.txt")

	if err != nil {
		assert.FailNow(t, err.Error())
		return
	}

	if _, err := file.Write([]byte("hello world")); err != nil {
		assert.FailNow(t, err.Error())
		return
	}

	if err := file.Close(); err != nil {
		assert.FailNow(t, err.Error())
	}
}

func teardown(t *testing.T) {
	if err := os.RemoveAll("./testdir"); err != nil {
		assert.Fail(t, err.Error())
	}
}
