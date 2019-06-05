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
		Action: withProjectDirectory(func(ctx *cli.Context, wd string) error {
			build, test, lint, err := command.Diff(wd)

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

			if len(lint) > 0 {
				fmt.Println("Modules to be linted:")
				for _, name := range lint {
					if _, err := fmt.Println(name); err != nil {
						return err
					}
				}
				fmt.Println()
			}

			return nil
		}),
	}
}
