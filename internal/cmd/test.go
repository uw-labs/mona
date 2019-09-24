package cmd

import (
	"github.com/davidsbond/mona/internal/command"
	"github.com/davidsbond/mona/internal/config"
	"github.com/davidsbond/mona/internal/deps"
	"github.com/urfave/cli"
)

// Test generates a command that runs tests for all apps with changes.
func Test() cli.Command {
	return cli.Command{
		Name:  "test",
		Usage: "Runs tests for all apps that have been created/changed since the last test run",
		Action: withModAndProject(func(ctx *cli.Context, mod deps.Module, pj *config.ProjectFile) error {
			return command.Test(mod, pj)
		}),
	}
}