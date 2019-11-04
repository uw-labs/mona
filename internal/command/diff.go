package command

import (
	"github.com/davidsbond/mona/internal/config"
)

// DiffSummary lists all the app names that changed.
type DiffSummary struct {
	All   []string
	Lint  []string
	Test  []string
	Build []string
}

// Diff outputs the names of all apps where changes are detected.
func Diff(pj *config.Project) (summary DiffSummary, err error) {
	_, changes, err := getLockAndChangedApps(pj)
	if err != nil {
		return summary, err
	}

	for _, change := range changes {
		summary.All = append(summary.All, change.app.Name)
		if change.changeTypes[changeTypeLint] {
			summary.Lint = append(summary.Lint, change.app.Name)
		}
		if change.changeTypes[changeTypeTest] {
			summary.Test = append(summary.Test, change.app.Name)
		}
		if change.changeTypes[changeTypeBuild] {
			summary.Build = append(summary.Build, change.app.Name)
		}
	}

	return summary, nil
}
