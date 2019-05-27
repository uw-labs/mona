package main

import (
	"errors"
	"fmt"
	"os"
	"sort"

	"github.com/davidsbond/mona/cmd"
	"github.com/davidsbond/mona/internal/files"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Usage = "A monorepo management tool"
	app.Author = "David Bond"
	app.Email = "davidsbond93@gmail.com"
	app.Copyright = "2019 David Bond"

	app.Commands = []cli.Command{
		cmd.Init(),
		cmd.AddModule(),
		cmd.Diff(),
		cmd.Build(),
		cmd.Test(),
		cmd.Lint(),
	}

	// Try to load the project file before every command to ensure we're
	// in a valid project directory
	app.Before = func(ctx *cli.Context) error {
		if _, err := files.LoadProjectFile(); err == files.ErrNoProject {
			app.CustomAppHelpTemplate = "No mona.yml file in current directory"
			return errors.New("")
		} else if err != nil {
			app.CustomAppHelpTemplate = err.Error()
			return errors.New("")
		}

		return nil
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	if err := app.Run(os.Args); err != nil {
		fmt.Println(err.Error())
	}
}
