// Package cmd contains definitions for executable commands and is responsible
// for the validation of flags and arguments.
package cmd

import (
	"github.com/davidsbond/mona/internal/config"
	"github.com/urfave/cli"
)

// The ActionFunc type is a method that takes a CLI context and the
// current project as an argument and returns a single error.
type ActionFunc func(ctx *cli.Context, p *config.ProjectFile) error

func withProject(fn ActionFunc) cli.ActionFunc {
	return func(ctx *cli.Context) error {
		_, project, err := getRootAndProject(ctx)
		if err != nil {
			return err
		}

		return fn(ctx, project)
	}
}

func getRootAndProject(ctx *cli.Context) (root string, project *config.ProjectFile, err error) {
	wd := ctx.GlobalString("wd")

	root, err = config.GetProjectRoot(wd)
	if err != nil {
		return "", nil, err
	}

	project, err = config.LoadProjectFile(root)

	return root, project, err
}
