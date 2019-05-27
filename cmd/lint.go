package cmd

import (
	"os"

	"github.com/davidsbond/mona/internal/command"
	"github.com/urfave/cli"
)

// Lint generates a command that lints all new/modified modules within the project.
func Lint() cli.Command {
	return cli.Command{
		Name:  "lint",
		Usage: "Lints any new/modified modules",
		Action: func(Ctx *cli.Context) error {
			wd, err := os.Getwd()

			if err != nil {
				return err
			}

			return command.Lint(wd)
		},
	}
}
