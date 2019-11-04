package main

import (
	"os"
	"sort"
	"strconv"
	"time"

	"github.com/apex/log"
	apexcli "github.com/apex/log/handlers/cli"
	"github.com/urfave/cli"

	"github.com/uw-labs/mona/internal/cmd"
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
	app.Usage = "A go monorepo management tool"
	app.Authors = []cli.Author{{
		Name:  "David Bond",
		Email: "davidsbond93@gmail.com",
	}, {
		Name: "Michal Bock",
	}}
	app.Copyright = "2019 David Bond"
	app.Version = version
	app.Compiled = time.Unix(compileTime, 0)

	app.Commands = []cli.Command{
		cmd.Init(),
		cmd.AddApp(),
		cmd.Diff(),
		cmd.Run(),
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
		cli.BoolFlag{
			Name:   "fail-fast",
			Usage:  "If set, causes a command to terminate after encountering first error",
			EnvVar: "MONA_FAIL_FAST",
		},
		cli.BoolFlag{
			Name:   "verbose",
			Usage:  "If set, enables verbose logging output",
			EnvVar: "MONA_VERBOSE",
		},
	}

	app.Before = func(ctx *cli.Context) error {
		log.SetHandler(apexcli.Default)

		if ctx.Bool("verbose") {
			log.SetLevel(log.DebugLevel)
		} else {
			log.SetLevel(log.InfoLevel)
		}

		log.Infof("%s v%s", app.Name, app.Version)

		wd, err := os.Getwd()

		if err != nil {
			return err
		}

		return ctx.Set("wd", wd)
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err.Error())
	}
}
