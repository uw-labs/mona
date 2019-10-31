package command

import (
	"github.com/davidsbond/mona/internal/config"
)

// Diff outputs the names of all apps where changes are detected.
func Diff(pj *config.Project) ([]string, []string, []string, error) {
	build, err := getChangedApps(pj, changeTypeBuild)

	if err != nil {
		return nil, nil, nil, err
	}

	test, err := getChangedApps(pj, changeTypeTest)

	if err != nil {
		return nil, nil, nil, err
	}

	lint, err := getChangedApps(pj, changeTypeLint)

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
