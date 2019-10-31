package cmd

import (
	"github.com/apex/log"
	"github.com/urfave/cli"

	"github.com/davidsbond/mona/internal/command"
)

// Diff generates a cli command that prints out apps that have changed within
// the current project.
func Diff() cli.Command {
	return cli.Command{
		Name:  "diff",
		Usage: "Outputs all apps where changes are detected",
		Action: withProject(func(ctx *cli.Context, cfg command.Config) error {
			lint, test, build, err := command.Diff(cfg.Project)

			if err != nil {
				return err
			}

			log.Infof("%d app(s) to be linted", len(lint))
			log.Infof("%d app(s) to be tested", len(test))
			log.Infof("%d app(s) to be built", len(build))

			return nil
		}),
	}
}
