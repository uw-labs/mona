package cmd

import (
	"runtime"

	"github.com/davidsbond/mona/internal/command"
	"github.com/davidsbond/mona/internal/files"
	"github.com/urfave/cli"
)

// Build generates a cli command that builds any modified modules within
// the project.
func Build() cli.Command {
	return cli.Command{
		Name:  "build",
		Usage: "Builds any modified modules within the project",
		Flags: []cli.Flag{
			cli.IntFlag{
				Name:  "parallelism",
				Usage: "Determines the number of threads to use to use when building modules",
				Value: runtime.NumCPU(),
			},
		},
		Action: withProject(func(ctx *cli.Context, pj *files.ProjectFile) error {
			return command.Build(pj, ctx.Int("parallelism"))
		}),
	}
}
