package cmd

import (
	"github.com/davidsbond/mona/internal/command"
	"github.com/davidsbond/mona/internal/deps"
	"github.com/davidsbond/mona/internal/files"
	"github.com/urfave/cli"
)

// Lint generates a command that lints all new/modified modules within the project.
func Lint() cli.Command {
	return cli.Command{
		Name:  "lint",
		Usage: "Lints any new/modified modules",
		Action: withModAndProject(func(ctx *cli.Context, mod deps.Module, pj *files.ProjectFile) error {
			return command.Lint(mod, pj)
		}),
	}
}
