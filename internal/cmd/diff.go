package cmd

import (
	"github.com/apex/log"

	"github.com/davidsbond/mona/internal/command"
	"github.com/davidsbond/mona/internal/config"
	"github.com/urfave/cli"
)

// Diff generates a cli command that prints out apps that have changed within
// the current project.
func Diff() cli.Command {
	return cli.Command{
		Name:  "diff",
		Usage: "Outputs all apps where changes are detected",
		Action: withProject(func(ctx *cli.Context, pj *config.ProjectFile) error {
			build, test, lint, err := command.Diff(pj)

			if err != nil {
				return err
			}

			log.Infof("%d app(s) to be built", len(build))
			log.Infof("%d app(s) to be tested", len(test))
			log.Infof("%d app(s) to be linted", len(lint))

			return nil
		}),
	}
}
