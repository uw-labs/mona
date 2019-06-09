package main

import (
	"os"
	"sort"
	"strconv"
	"time"

	"github.com/davidsbond/mona/cmd"
	"github.com/davidsbond/mona/internal/output"
	"github.com/urfave/cli"
)

var (
	version     string
	compiled    string
	compileTime int64
)

func init() {
	compileTime, _ = strconv.ParseInt(compiled, 10, 64)
}

func main() {
	app := cli.NewApp()
	app.Usage = "A monorepo management tool"
	app.Author = "David Bond"
	app.Email = "davidsbond93@gmail.com"
	app.Copyright = "2019 David Bond"
	app.Version = version
	app.Compiled = time.Unix(compileTime, 0)

	app.Commands = []cli.Command{
		cmd.Init(),
		cmd.AddModule(),
		cmd.Diff(),
		cmd.Build(),
		cmd.Test(),
		cmd.Lint(),
	}

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "wd",
			Hidden: true,
			Usage:  "Flag that contains the current working directory",
		},
	}

	app.Before = func(ctx *cli.Context) error {
		wd, err := os.Getwd()

		if err != nil {
			return err
		}

		return ctx.Set("wd", wd)
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	if err := app.Run(os.Args); err != nil {
		if err = output.WriteError(os.Stdout, err); err != nil {
			panic(err)
		}
	}
}
