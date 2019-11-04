// Package cmd contains definitions for executable commands and is responsible
// for the validation of flags and arguments.
package cmd

import (
	"github.com/apex/log"
	"github.com/urfave/cli"
	"github.com/uw-labs/mona/internal/command"
	"github.com/uw-labs/mona/internal/config"
)

// The ActionFunc type is a method that takes a CLI context and the
// current project as an argument and returns a single error.
type ActionFunc func(ctx *cli.Context, cfg command.Config) error

func withCMDConfig(fn ActionFunc) cli.ActionFunc {
	return func(ctx *cli.Context) (err error) {
		wd := ctx.GlobalString("wd")

		branch := ctx.GlobalString("compare-git-branch")
		log.Debugf("Comparing against %s branch.", branch)

		root, err := config.GetProjectRoot(wd)
		if err != nil {
			return err
		}

		cfg, err := command.NewConfig(root, branch)
		if err != nil {
			return err
		}
		cfg.FailFast = ctx.GlobalBool("fail-fast")

		return fn(ctx, cfg)
	}
}
