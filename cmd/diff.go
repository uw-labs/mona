package cmd

import (
	"fmt"

	"github.com/davidsbond/mona/internal/command"

	"github.com/urfave/cli"
)

// Diff generates a cli command that prints out modules that have changed within
// the current project.
func Diff() cli.Command {
	return cli.Command{
		Name:  "diff",
		Usage: "Outputs all modules where changes are detected",
		Action: func(ctx *cli.Context) error {
			build, test, err := command.Diff()

			if err != nil {
				return err
			}

			if len(build) > 0 {
				fmt.Println("Modules to be built:")
				for _, name := range build {
					if _, err := fmt.Println(name); err != nil {
						return err
					}
				}
				fmt.Println()
			}

			if len(test) > 0 {
				fmt.Println("Modules to be tested:")
				for _, name := range test {
					if _, err := fmt.Println(name); err != nil {
						return err
					}
				}
				fmt.Println()
			}

			return nil
		},
	}
}
