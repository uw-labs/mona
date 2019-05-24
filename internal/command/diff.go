package command

// Diff outputs the names of all modules where changes are detected.
func Diff() ([]string, error) {
	changed, err := getChangedModules()

	if err != nil {
		return nil, err
	}

	var names []string
	for _, module := range changed {
		names = append(names, module.Name)
	}

	return names, nil
}
