package cmd

import (
	"os"

	"github.com/davidsbond/mona/internal/command"
	"github.com/urfave/cli"
)

// Test generates a command that runs tests for all modules with changes.
func Test() cli.Command {
	return cli.Command{
		Name:  "test",
		Usage: "Runs tests for all modules that have been created/changed since the last test run",
		Action: func(ctx *cli.Context) error {
			wd, err := os.Getwd()

			if err != nil {
				return err
			}

			return command.Test(wd)
		},
	}
}
