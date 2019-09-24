package command

import (
	"github.com/davidsbond/mona/internal/config"
	"github.com/davidsbond/mona/internal/deps"
)

// Diff outputs the names of all apps where changes are detected.
func Diff(mod deps.Module, pj *config.ProjectFile) ([]string, []string, []string, error) {
	build, err := getChangedApps(mod, pj, changeTypeBuild)

	if err != nil {
		return nil, nil, nil, err
	}

	test, err := getChangedApps(mod, pj, changeTypeTest)

	if err != nil {
		return nil, nil, nil, err
	}

	lint, err := getChangedApps(mod, pj, changeTypeLint)

	if err != nil {
		return nil, nil, nil, err
	}

	var buildNames []string
	for _, app := range build {
		buildNames = append(buildNames, app.Name)
	}

	var testNames []string
	for _, app := range test {
		testNames = append(testNames, app.Name)
	}

	var lintNames []string
	for _, app := range lint {
		lintNames = append(lintNames, app.Name)
	}

	return buildNames, testNames, lintNames, nil
}
