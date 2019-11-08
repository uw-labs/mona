package main

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"time"

	"github.com/apex/log"
	"github.com/urfave/cli"

	"github.com/uw-labs/mona/internal/git"
	"github.com/uw-labs/mona/internal/golang"
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

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "compare-git-branch",
			EnvVar: "MONA_COMPARE_GIT_BRANCH",
			Value:  "master",
		},
	}
	app.Commands = []cli.Command{{
		Name: "diff",
		Subcommands: []cli.Command{{
			Name: "test",
			Action: withModAndDiff(func(ctx *cli.Context, mod golang.Module, diff git.GoDiff) error {
				for pkg := range diff.Packages {
					fmt.Println(pkg)
				}
				cmds, err := golang.FindAllCommands(mod.Name)
				if err != nil {
					return err
				}
				// Find all commands that depend on the changed packages.
				for _, cmd := range cmds {
					changed, err := diff.Changed(cmd, mod)
					if err != nil {
						return err
					}
					if changed && !diff.Packages[cmd] {
						fmt.Println(cmd)
					}
				}
				return nil
			}),
		}, {
			Name: "build",
			Action: withModAndDiff(func(ctx *cli.Context, mod golang.Module, diff git.GoDiff) error {
				cmds, err := golang.FindAllCommands(mod.Name)
				if err != nil {
					return err
				}
				// Find all commands that depend on the changed packages.
				for _, cmd := range cmds {
					changed, err := diff.Changed(cmd, mod)
					if err != nil {
						return err
					}
					if changed {
						fmt.Println(cmd)
					}
				}
				return nil
			}),
		}},
	}}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err.Error())
	}
}

func withModAndDiff(f func(ctx *cli.Context, mod golang.Module, diff git.GoDiff) error) func(*cli.Context) error {
	return func(ctx *cli.Context) error {
		mod, err := golang.ParseModuleFile("go.mod")
		if err != nil {
			return err
		}
		diff, err := git.GetGoDiff(mod, ctx.GlobalString("compare-git-branch"))
		if err != nil {
			return err
		}

		return f(ctx, mod, diff)
	}
}
