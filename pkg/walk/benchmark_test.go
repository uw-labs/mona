package walk_test

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/davidsbond/mona/pkg/walk"
	"github.com/stretchr/testify/assert"
)

func BenchmarkFast(b *testing.B) {
	b.StopTimer()
	setupBenchmark(b, "./bench")

	for i := 0; i < b.N; i++ {
		b.StartTimer()

		err := walk.Fast("./bench", func(path string, info os.FileInfo, err error) error {
			return nil
		}, runtime.NumCPU())

		b.StopTimer()

		if err != nil {
			assert.FailNow(b, err.Error())
			return
		}
	}

	b.StopTimer()
	if err := os.RemoveAll("./bench"); err != nil {
		assert.Fail(b, err.Error())
	}
}

func setupBenchmark(b *testing.B, dir string) {
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		assert.FailNow(b, err.Error())
		return
	}

	for i := 0; i < b.N; i++ {
		path := filepath.Join(dir, fmt.Sprintf("%v.txt", i))

		file, err := os.Create(path)

		if err != nil {
			assert.FailNow(b, err.Error())
			return
		}

		if err := file.Close(); err != nil {
			assert.FailNow(b, err.Error())
			return
		}
	}
}
