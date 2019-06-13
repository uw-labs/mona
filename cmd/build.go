package cmd

import (
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
		Action: withProject(func(ctx *cli.Context, pj *files.ProjectFile) error {
			return command.Build(pj)
		}),
	}
}
