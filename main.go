package main

import (
	"github.com/davidsbond/mona/cmd"
	"os"
	"sort"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Usage = "A monorepo management tool"
	app.Commands = []cli.Command{
		cmd.Init(),
		cmd.AddModule(),
		cmd.Diff(),
		cmd.Build(),
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	if err := app.Run(os.Args); err != nil {
		panic(err)
	}
}
