package cmd

import (
	"github.com/davidsbond/mona/internal/command"
	"github.com/urfave/cli"
)

// Lint generates a command that lints all new/modified modules within the project.
func Lint() cli.Command {
	return cli.Command{
		Name:  "lint",
		Usage: "Lints any new/modified modules",
		Action: withProjectDirectory(func(ctx *cli.Context, wd string) error {
			return command.Lint(wd)
		}),
	}
}
