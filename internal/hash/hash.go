package hash

import (
	"crypto/md5"
	"encoding/base64"

	"github.com/uw-labs/mona/internal/app"
	"github.com/uw-labs/mona/pkg/hashdir"
)

func GetForApp(appInfo *app.App, excludes ...string) (string, error) {
	hash := md5.New()
	for _, dep := range appInfo.Deps.External {
		if _, err := hash.Write([]byte(dep)); err != nil {
			return "", err
		}
	}

	for _, dep := range appInfo.Deps.Internal {
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

	appHash, err := hashdir.Generate(appInfo.Location, excludes...)
	if err != nil {
		return "", err
	}
	if _, err = hash.Write(appHash); err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(hash.Sum(nil)), nil
}
