package cmd

import (
	"errors"
	"path/filepath"
	"sort"

	"github.com/davidsbond/mona/internal/files"

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
				return errors.New("argument 'location' is required")
			}

			if ctx.String("name") != "" {
				return nil
			}

			return ctx.Set("name", filepath.Base(dir))
		},
		Action: withProject(func(ctx *cli.Context, pj *files.ProjectFile) error {
			return command.AddModule(
				pj,
				ctx.String("name"),
				ctx.Args().First())
		}),
	}

	sort.Sort(cli.FlagsByName(cmd.Flags))
	return cmd
}
