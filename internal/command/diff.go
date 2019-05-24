package command

import (
	"fmt"
)

// Diff outputs the names of all modules where changes are detected.
func Diff() error {
	changed, err := getChangedModules()

	if err != nil {
		return err
	}

	for _, module := range changed {
		if _, err := fmt.Println(module.Name); err != nil {
			return err
		}
	}

	return nil
}
