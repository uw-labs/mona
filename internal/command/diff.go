package command

import "github.com/davidsbond/mona/internal/files"

// Diff outputs the names of all modules where changes are detected.
func Diff(pj *files.ProjectFile, parallelism int) ([]string, []string, []string, error) {
	build, err := getChangedModules(pj, changeTypeBuild, parallelism)

	if err != nil {
		return nil, nil, nil, err
	}

	test, err := getChangedModules(pj, changeTypeTest, parallelism)

	if err != nil {
		return nil, nil, nil, err
	}

	lint, err := getChangedModules(pj, changeTypeLint, parallelism)

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
