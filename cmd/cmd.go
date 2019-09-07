// Package cmd contains definitions for executable commands and is responsible
// for the validation of flags and arguments.
package cmd

import (
	"path/filepath"

	"github.com/davidsbond/mona/internal/deps"
	"github.com/davidsbond/mona/internal/files"
	"github.com/urfave/cli"
)

type (
	// The ActionFunc type is a method that takes a CLI context and the
	// current project as an argument and returns a single error.
	ActionFunc func(ctx *cli.Context, p *files.ProjectFile) error
	// The BuildActionFunc type is same as ActionFunc, but it also takes
	// go module config as an argument.
	BuildActionFunc func(ctx *cli.Context, mod deps.Module, p *files.ProjectFile) error
)

func withProject(fn ActionFunc) cli.ActionFunc {
	return func(ctx *cli.Context) error {
		_, project, err := getRootAndProject(ctx)
		if err != nil {
			return err
		}

		return fn(ctx, project)
	}
}

func withModAndProject(fn BuildActionFunc) cli.ActionFunc {
	return func(ctx *cli.Context) error {
		root, project, err := getRootAndProject(ctx)
		if err != nil {
			return err
		}

		mod, err := deps.ParseModule(filepath.Join(root, "go.mod"))
		if err != nil {
			return err
		}

		return fn(ctx, mod, project)
	}
}

func getRootAndProject(ctx *cli.Context) (root string, project *files.ProjectFile, err error) {
	wd := ctx.GlobalString("wd")

	root, err = files.GetProjectRoot(wd)
	if err != nil {
		return "", nil, err
	}

	project, err = files.LoadProjectFile(root)

	return root, project, err
}
