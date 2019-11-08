package cmd

import (
	"github.com/apex/log"
	"github.com/urfave/cli"

	"github.com/uw-labs/mona/internal/command"
)

// Diff generates a cli command that prints out apps that have changed within
// the current project.
func Diff() cli.Command {
	return cli.Command{
		Name:  "diff",
		Usage: "Outputs all apps where changes are detected",
		Action: withCMDConfig(func(ctx *cli.Context, cfg command.Config) error {
			summary, err := command.Diff(cfg)

			if err != nil {
				return err
			}

			log.Infof("%d app(s) in project", len(summary.All))
			for _, appInfo := range cfg.Apps {
				log.Debugf("\t%s", appInfo.Name)
				for _, dep := range appInfo.Deps.Internal {
					log.Debugf("\t\t%s", dep)
				}
			}
			log.Infof("%d app(s) changed", len(summary.ChangedApps))
			for _, appName := range summary.ChangedApps {
				log.Debugf("\t%s", appName)
			}
			log.Infof("%d pkg(s) changed", len(cfg.Diff.Packages))
			for pkg := range cfg.Diff.Packages {
				log.Debugf("\t%s", pkg)
			}

			return nil
		}),
	}
}
