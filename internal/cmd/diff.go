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
			summary, err := command.Diff(cfg.Project)

			if err != nil {
				return err
			}

			log.Infof("%d app(s) in project", len(summary.All))
			log.Infof("%d app(s) to be linted", len(summary.Lint))
			log.Infof("%d app(s) to be tested", len(summary.Test))
			log.Infof("%d app(s) to be built", len(summary.Build))

			return nil
		}),
	}
}
