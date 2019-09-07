package command

import (
	"github.com/davidsbond/mona/internal/deps"
	"github.com/davidsbond/mona/internal/files"
)

// Diff outputs the names of all modules where changes are detected.
func Diff(mod deps.Module, pj *files.ProjectFile) ([]string, []string, []string, error) {
	build, err := getChangedModules(mod, pj, changeTypeBuild)

	if err != nil {
		return nil, nil, nil, err
	}

	test, err := getChangedModules(mod, pj, changeTypeTest)

	if err != nil {
		return nil, nil, nil, err
	}

	lint, err := getChangedModules(mod, pj, changeTypeLint)

	if err != nil {
		return nil, nil, nil, err
	}

	var buildNames []string
	for _, module := range build {
		buildNames = append(buildNames, module.Name)
	}

	var testNames []string
	for _, module := range test {
		testNames = append(testNames, module.Name)
	}

	var lintNames []string
	for _, module := range lint {
		lintNames = append(lintNames, module.Name)
	}

	return buildNames, testNames, lintNames, nil
}
