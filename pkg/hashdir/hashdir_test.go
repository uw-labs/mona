package hashdir_test

import (
	"encoding/base64"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/uw-labs/mona/pkg/hashdir"
)

func TestGenerate(t *testing.T) {
	tt := []struct {
		Name      string
		Directory string
		Expected  string
	}{
		{
			Name:      "It should generate a hash of a directory",
			Directory: "./testdir",
			Expected:  "XrY7u+Ae7tCTyyK7j1rNww==",
		},
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			setup(t)
			defer teardown(t)

			hash, err := hashdir.Generate(tc.Directory)

			if err != nil {
				assert.Fail(t, err.Error())
				return
			}

			assert.Equal(t, tc.Expected, base64.StdEncoding.EncodeToString(hash))
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
