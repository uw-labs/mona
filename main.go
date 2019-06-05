package main

import (
	"os"
	"sort"

	"github.com/davidsbond/mona/cmd"
	"github.com/davidsbond/mona/internal/output"
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

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	if err := app.Run(os.Args); err != nil {
		output.WriteError(os.Stdout, err)
	}
}
