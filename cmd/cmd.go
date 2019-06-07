// Package cmd contains definitions for executable commands and is responsible
// for the validation of flags and arguments.
package cmd

import (
	"os"

	"github.com/davidsbond/mona/internal/files"
	"github.com/urfave/cli"
)

type (
	// The ActionFunc type is a method that takes a CLI context and the
	// current project as an argument and returns a single error.
	ActionFunc func(ctx *cli.Context, p *files.ProjectFile) error
)

func withProject(fn ActionFunc) cli.ActionFunc {
	return func(ctx *cli.Context) error {
		wd, err := os.Getwd()

		if err != nil {
			return err
		}

		root, err := files.GetProjectRoot(wd)

		if err != nil {
			return err
		}

		project, err := files.LoadProjectFile(root)

		if err != nil {
			return err
		}

		return fn(ctx, project)
	}
}
