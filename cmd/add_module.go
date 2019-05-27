package cmd

import (
	"errors"
	"os"
	"path/filepath"
	"sort"

	"github.com/davidsbond/mona/internal/command"
	"github.com/urfave/cli"
)

// AddModule generates a cli command that creates new mona modules within a project.
func AddModule() cli.Command {
	cmd := cli.Command{
		Name:      "add-module",
		Usage:     "Initializes a new module at the provided path",
		ArgsUsage: "<location>",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "name",
				Usage: "The name of the module, defaults to base of the module directory",
			},
		},
		// Before the command is called, check if the name provided is blank. If so, set it
		// to the base of the given location argument.
		Before: func(ctx *cli.Context) error {
			dir := ctx.Args().First()

			if dir == "" {
				return errors.New("Argument 'location' is required")
			}

			if ctx.String("name") != "" {
				return nil
			}

			return ctx.Set("name", filepath.Base(dir))
		},
		Action: func(ctx *cli.Context) error {
			wd, err := os.Getwd()

			if err != nil {
				return err
			}

			location := ctx.Args().First()
			name := ctx.String("name")

			return command.AddModule(wd, name, location)
		},
	}

	sort.Sort(cli.FlagsByName(cmd.Flags))
	return cmd
}
