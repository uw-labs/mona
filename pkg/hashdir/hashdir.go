// Package hashdir contains methods for generating hashes of directories
package hashdir

import (
	"crypto/md5"
	"encoding/base64"
	"io"
	"os"
	"path/filepath"
	"sync"

	"github.com/davidsbond/mona/pkg/walk"
)

// Generate creates a new hash for a given path. The path is walked and a hash is
// created based on all files found in the path. If a file matches one specified in
// the 'excludes' parameter it is not used to generate the hash.
func Generate(location string, parallelism int, excludes ...string) (string, error) {
	var mux sync.Mutex
	hash := md5.New()

	err := walk.Fast(location, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		for _, exclude := range excludes {
			ok, err := filepath.Match(exclude, info.Name())

			if err != nil {
				return err
			}

			if ok {
				return nil
			}
		}

		if info.IsDir() {
			return nil
		}

		file, err := os.Open(path)

		if err != nil {
			return err
		}

		mux.Lock()
		defer mux.Unlock()
		if _, err := io.Copy(hash, file); err != nil {
			return err
		}

		return file.Close()
	}, parallelism)

	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(hash.Sum(nil)), nil
}
