package cmd

import (
	"github.com/davidsbond/mona/internal/command"
	"github.com/urfave/cli"
)

// Test generates a command that runs tests for all modules with changes.
func Test() cli.Command {
	return cli.Command{
		Name:  "test",
		Usage: "Runs tests for all modules that have been created/changed since the last test run",
		Action: func(ctx *cli.Context) error {
			return command.Test()
		},
	}
}
