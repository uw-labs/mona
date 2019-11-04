// Package cmd contains definitions for executable commands and is responsible
// for the validation of flags and arguments.
package cmd

import (
	"github.com/urfave/cli"

	"github.com/uw-labs/mona/internal/command"
	"github.com/uw-labs/mona/internal/config"
)

// The ActionFunc type is a method that takes a CLI context and the
// current project as an argument and returns a single error.
type ActionFunc func(ctx *cli.Context, cfg command.Config) error

func withProject(fn ActionFunc) cli.ActionFunc {
	return func(ctx *cli.Context) error {
		wd := ctx.GlobalString("wd")

		root, err := config.GetProjectRoot(wd)
		if err != nil {
			return err
		}

		project, err := config.LoadProject(root)
		if err != nil {
			return err
		}

		return fn(ctx, command.Config{
			Project:  project,
			FailFast: ctx.GlobalBool("fail-fast"),
		})
	}
}
