package cmd

import (
	"github.com/apex/log"

	"github.com/davidsbond/mona/internal/command"
	"github.com/davidsbond/mona/internal/deps"
	"github.com/davidsbond/mona/internal/files"
	"github.com/urfave/cli"
)

// Diff generates a cli command that prints out modules that have changed within
// the current project.
func Diff() cli.Command {
	return cli.Command{
		Name:  "diff",
		Usage: "Outputs all modules where changes are detected",
		Action: withModAndProject(func(ctx *cli.Context, mod deps.Module, pj *files.ProjectFile) error {
			build, test, lint, err := command.Diff(mod, pj)

			if err != nil {
				return err
			}

			log.Infof("%d module(s) to be built", len(build))
			log.Infof("%d module(s) to be tested", len(test))
			log.Infof("%d module(s) to be linted", len(lint))

			return nil
		}),
	}
}
