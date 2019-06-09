package cmd

import (
	"runtime"

	"github.com/davidsbond/mona/internal/command"
	"github.com/davidsbond/mona/internal/files"
	"github.com/urfave/cli"
)

// Lint generates a command that lints all new/modified modules within the project.
func Lint() cli.Command {
	return cli.Command{
		Name:  "lint",
		Usage: "Lints any new/modified modules",
		Flags: []cli.Flag{
			cli.IntFlag{
				Name:  "parallelism",
				Usage: "Determines the number of threads to use when linting modules",
				Value: runtime.NumCPU(),
			},
		},
		Action: withProject(func(ctx *cli.Context, pj *files.ProjectFile) error {
			return command.Lint(pj, ctx.Int("parallelism"))
		}),
	}
}
