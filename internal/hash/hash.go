package hash

import (
	"crypto/md5"
	"encoding/base64"
	"sort"

	"github.com/davidsbond/mona/internal/deps"
	"github.com/davidsbond/mona/pkg/hashdir"
)

func GetAppDeps(mod deps.Module, appPath string, excludes ...string) (string, error) {
	appDeps, err := deps.GetAppDeps(mod, appPath)
	if err != nil {
		return "", err
	}
	sort.Sort(sort.StringSlice(appDeps.Internal))
	sort.Sort(sort.StringSlice(appDeps.External))

	hash := md5.New()
	for _, dep := range appDeps.External {
		if _, err := hash.Write([]byte(dep)); err != nil {
			return "", err
		}
	}

	for _, dep := range appDeps.Internal {
		depHash, err := hashdir.Generate(dep, excludes...)
		if err != nil {
			return "", err
		}
		if _, err := hash.Write(depHash); err != nil {
			return "", err
		}
		// Exclude current dependency going forward to avoid
		// rehashing the same directory over and over.
		excludes = append(excludes, dep)
	}

	appHash, err := hashdir.Generate(appPath, excludes...)
	if err != nil {
		return "", err
	}
	if _, err = hash.Write(appHash); err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(hash.Sum(nil)), nil
}
