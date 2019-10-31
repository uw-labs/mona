package cmd

import (
	"errors"
	"path/filepath"
	"sort"

	"github.com/apex/log"
	"github.com/urfave/cli"

	"github.com/davidsbond/mona/internal/command"
)

// AddApp generates a cli command that creates new mona apps within a project.
func AddApp() cli.Command {
	cmd := cli.Command{
		Name:      "add-app",
		Usage:     "Initializes a new app at the provided path",
		ArgsUsage: "<location>",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "name",
				Usage: "The name of the app, defaults to base of the app directory",
			},
		},
		// Before the command is called, check if the name provided is blank. If so, set it
		// to the base of the given location argument.
		Before: func(ctx *cli.Context) error {
			dir := ctx.Args().First()

			if dir == "" {
				return errors.New("argument 'location' is required")
			}

			if ctx.String("name") != "" {
				return nil
			}

			return ctx.Set("name", filepath.Base(dir))
		},
		Action: withProject(func(ctx *cli.Context, cfg command.Config) error {
			name := ctx.String("name")
			dir := ctx.Args().First()

			if err := command.AddApp(cfg.Project, name, dir); err != nil {
				return err
			}

			log.Infof("Created new app %s at %s", name, dir)
			return nil
		}),
	}

	sort.Sort(cli.FlagsByName(cmd.Flags))
	return cmd
}
