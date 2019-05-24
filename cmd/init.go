package cmd

import (
	"github.com/davidsbond/mona/internal/command"
	"os"
	"path/filepath"

	"github.com/urfave/cli"
)

const (
	defaultVersion = "v1"
)

// Init generates a cli command for initializing new mona projects.
func Init() cli.Command {
	return cli.Command{
		Name:  "init",
		Usage: "Initializes a new mona project in the current working directory",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "version",
				Usage: "The mona version used for this project",
				Value: defaultVersion,
			},
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

			wd, err := os.Getwd()

			if err != nil {
				return err
			}

			return ctx.Set("name", filepath.Base(wd))
		},
		Action: func(ctx *cli.Context) error {
			name := ctx.String("name")
			version := ctx.String("version")

			return command.Init(name, version)
		},
	}
}
