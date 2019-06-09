package cmd

import (
	"runtime"

	"github.com/davidsbond/mona/internal/command"
	"github.com/davidsbond/mona/internal/files"
	"github.com/urfave/cli"
)

// Test generates a command that runs tests for all modules with changes.
func Test() cli.Command {
	return cli.Command{
		Name:  "test",
		Usage: "Runs tests for all modules that have been created/changed since the last test run",
		Flags: []cli.Flag{
			cli.IntFlag{
				Name:  "parallelism",
				Usage: "Determines the number of threads to use when testing modules",
				Value: runtime.NumCPU(),
			},
		},
		Action: withProject(func(ctx *cli.Context, pj *files.ProjectFile) error {
			return command.Test(pj, ctx.Int("parallelism"))
		}),
	}
}
