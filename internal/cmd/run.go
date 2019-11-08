package cmd

import (
	"github.com/urfave/cli"

	"github.com/uw-labs/mona/internal/command"
)

// Run generates a cli command that lints, tests and builds any modified apps
// within the project.
func Run() cli.Command {
	return cli.Command{
		Name:  "run",
		Usage: "Run lints, tests and builds any modified apps within the project",
		Action: withCMDConfig(func(ctx *cli.Context, cfg command.Config) error {
			return command.Run(cfg)
		}),
	}
}

// Lint generates a command that lints all new/modified apps within the project.
func Lint() cli.Command {
	return cli.Command{
		Name:  "lint",
		Usage: "Lints any new/modified apps",
		Action: withCMDConfig(func(ctx *cli.Context, cfg command.Config) error {
			return command.Lint(cfg)
		}),
	}
}

// Test generates a command that runs tests for all apps with changes.
func Test() cli.Command {
	return cli.Command{
		Name:  "test",
		Usage: "Runs tests for all apps that have been created/changed since the last test run",
		Action: withCMDConfig(func(ctx *cli.Context, cfg command.Config) error {
			return command.Test(cfg)
		}),
	}
}

// Build generates a cli command that builds any modified apps within
// the project.
func Build() cli.Command {
	return cli.Command{
		Name:  "build",
		Usage: "Builds any modified apps within the project",
		Action: withCMDConfig(func(ctx *cli.Context, cfg command.Config) error {
			return command.Build(cfg)
		}),
	}
}
