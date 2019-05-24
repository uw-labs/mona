// Package hash contains methods for generating hashes
package hash

import (
	"crypto/md5"
	"encoding/base64"
	"io"
	"os"
	"path/filepath"
)

// Generate creates a new hash for a given path. The path is walked and a hash is
// created based on all files found in the path. If a file matches one specified in
// the 'exludes' parameter it is not used to generate the hash.
func Generate(location string, excludes ...string) (string, error) {
	hash := md5.New()

	err := filepath.Walk(location, func(path string, info os.FileInfo, err error) error {
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

		defer file.Close()
		if _, err := io.Copy(hash, file); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(hash.Sum(nil)), nil
}
