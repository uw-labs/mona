package cmd

import (
	"os"

	"github.com/davidsbond/mona/internal/files"

	"github.com/davidsbond/mona/internal/output"

	"github.com/davidsbond/mona/internal/command"
	"github.com/urfave/cli"
)

// Diff generates a cli command that prints out modules that have changed within
// the current project.
func Diff() cli.Command {
	return cli.Command{
		Name:  "diff",
		Usage: "Outputs all modules where changes are detected",
		Action: withProject(func(ctx *cli.Context, pj *files.ProjectFile) error {
			build, test, lint, err := command.Diff(pj)

			if err != nil {
				return err
			}

			if err := output.WriteList(os.Stdout, "Modules to be built:", build); err != nil {
				return err
			}

			if err := output.WriteList(os.Stdout, "Modules to be tested:", test); err != nil {
				return err
			}

			if err := output.WriteList(os.Stdout, "Modules to be linted:", lint); err != nil {
				return err
			}

			return nil
		}),
	}
}
