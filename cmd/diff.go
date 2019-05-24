package cmd

import (
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
			return command.Diff()
		},
	}
}
