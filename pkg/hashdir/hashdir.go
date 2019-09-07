// Package hashdir contains methods for generating hashes of directories
package hashdir

import (
	"crypto/md5"
	"io"
	"os"
	"path/filepath"

	"github.com/gobwas/glob"
)

// GenerateString creates a new hash for a given path. The path is walked and a hash is
// created based on all files found in the path. If a file matches one specified in
// the 'excludes' parameter it is not used to generate the hash.
func Generate(location string, excludes ...string) ([]byte, error) {
	hash := md5.New()
	globs := make(map[string]glob.Glob)

	err := filepath.Walk(location, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		for _, exclude := range excludes {
			gl, ok := globs[exclude]

			if !ok {
				gl, err = glob.Compile(exclude, os.PathSeparator)

				if err != nil {
					return err
				}

				globs[exclude] = gl
			}

			if gl.Match(info.Name()) || gl.Match(path) {
				if info.IsDir() {
					return filepath.SkipDir
				}

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

		if _, err := io.Copy(hash, file); err != nil {
			return err
		}

		return file.Close()
	})

	if err != nil {
		return nil, err
	}

	return hash.Sum(nil), nil
}
