package command

// Diff outputs the names of all modules where changes are detected.
func Diff() ([]string, []string, error) {
	build, err := getChangedModules(changeTypeBuild)

	if err != nil {
		return nil, nil, err
	}

	test, err := getChangedModules(changeTypeTest)

	if err != nil {
		return nil, nil, err
	}

	var buildNames []string
	for _, module := range build {
		buildNames = append(buildNames, module.Name)
	}

	var testNames []string
	for _, module := range test {
		testNames = append(testNames, module.Name)
	}

	return buildNames, testNames, nil
}
