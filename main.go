package main

import (
	"fmt"
	"os"
	"sort"

	"github.com/davidsbond/mona/cmd"
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
		fmt.Println(err.Error())
	}
}
