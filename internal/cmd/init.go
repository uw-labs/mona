package cmd

import (
	"path/filepath"
	"sort"

	"github.com/apex/log"
	"github.com/urfave/cli"

	"github.com/uw-labs/mona/internal/command"
)

// Init generates a cli command for initializing new mona projects.
func Init() cli.Command {
	cmd := cli.Command{
		Name:  "init",
		Usage: "Initializes a new mona project in the current working directory",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "name",
				Usage: "The name of the mona project, defaults to base of current working directory",
			},
		},
		// Before the command is called, check if the name provided is blank. If so, set it
		// to the base of the current working directory.
		Before: func(ctx *cli.Context) error {
			if ctx.String("name") != "" {
				return nil
			}

			name := filepath.Base(ctx.GlobalString("wd"))
			return ctx.Set("name", name)
		},
		Action: func(ctx *cli.Context) error {
			name := ctx.String("name")
			wd := ctx.GlobalString("wd")

			if err := command.Init(wd, name); err != nil {
				return err
			}

			log.Infof("Initialized new project at %s", wd)
			return nil
		},
	}

	sort.Sort(cli.FlagsByName(cmd.Flags))
	return cmd
}
