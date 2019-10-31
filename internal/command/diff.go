package command

import (
	"github.com/davidsbond/mona/internal/config"
)

// Diff outputs the names of all apps where changes are detected.
func Diff(pj *config.Project) (lint []string, test []string, build []string, err error) {
	_, changes, err := getLockAndChangedApps(pj)
	if err != nil {
		return nil, nil, nil, err
	}

	for _, change := range changes {
		if change.changeTypes[changeTypeLint] {
			lint = append(lint, change.app.Name)
		}
		if change.changeTypes[changeTypeTest] {
			test = append(test, change.app.Name)
		}
		if change.changeTypes[changeTypeBuild] {
			build = append(build, change.app.Name)
		}
	}

	return lint, test, build, nil
}
