package cmd

import (
	"github.com/davidsbond/mona/internal/command"
	"github.com/davidsbond/mona/internal/config"
	"github.com/davidsbond/mona/internal/deps"
	"github.com/urfave/cli"
)

// Lint generates a command that lints all new/modified apps within the project.
func Lint() cli.Command {
	return cli.Command{
		Name:  "lint",
		Usage: "Lints any new/modified apps",
		Action: withModAndProject(func(ctx *cli.Context, mod deps.Module, pj *config.ProjectFile) error {
			return command.Lint(mod, pj)
		}),
	}
}
