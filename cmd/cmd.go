// Package cmd contains definitions for executable commands and is responsible
// for the validation of flags and arguments.
package cmd

import (
	"os"

	"github.com/davidsbond/mona/internal/files"
	"github.com/urfave/cli"
)

type (
	ActionFunc func(ctx *cli.Context, pd string) error
)

func withProjectDirectory(fn ActionFunc) cli.ActionFunc {
	return func(ctx *cli.Context) error {
		wd, err := os.Getwd()

		if err != nil {
			return err
		}

		root, err := files.GetProjectRoot(wd)

		if err != nil {
			return err
		}

		return fn(ctx, root)
	}
}
