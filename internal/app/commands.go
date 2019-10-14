package app

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/apex/log"
)

func (app *App) Build(binDir string) error {
	flags := append([]string{"build"}, app.Commands.Build.Flags...)
	flags = append(flags, fmt.Sprintf("-o=%s/%s", binDir, app.Name), app.Location)

	cmd := exec.Command("go", flags...)
	cmd.Env = append(os.Environ(), app.Commands.Build.EnvToList()...)
	log.Info(cmd.String())

	return executeCommand(cmd)
}

func (app *App) Lint(checked map[string]bool) error {
	flags := append([]string{"run"}, app.Commands.Lint.Flags...)
	flags = append(flags, app.Location)

	for _, dep := range app.Deps.Internal {
		if !checked[dep] {
			flags = append(flags, "./"+dep)
			checked[dep] = true
		}
	}

	cmd := exec.Command("golangci-lint", flags...)
	cmd.Env = append(os.Environ(), app.Commands.Lint.EnvToList()...)
	log.Info(cmd.String())

	return executeCommand(cmd)
}

func (app *App) Test(checked map[string]bool) error {
	flags := append([]string{"test"}, app.Commands.Test.Flags...)
	flags = append(flags, app.Location)

	for _, dep := range app.Deps.Internal {
		if !checked[dep] {
			flags = append(flags, "./"+dep)
			checked[dep] = true
		}
	}

	cmd := exec.Command("go", flags...)
	cmd.Env = append(os.Environ(), app.Commands.Test.EnvToList()...)
	log.Info(cmd.String())

	return executeCommand(cmd)
}
