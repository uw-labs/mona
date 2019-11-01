package app

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/apex/log"
)

func (app *App) Build(binDir string) error {
	flags := append([]string{"build"}, fmt.Sprintf("-o=%s/%s", binDir, app.Name))
	flags = append(flags, app.Commands.Build.AllFlags()...)
	flags = append(flags, app.Location)

	cmd := exec.Command("go", flags...)
	cmd.Env = append(os.Environ(), app.Commands.Build.EnvToList()...)
	log.Info(cmd.String())

	return executeCommand(cmd)
}

func (app *App) Lint(checked map[string]bool) error {
	flags := append([]string{"run"}, app.Commands.Lint.Flags...)
	flags = append(flags, app.Location)

	checking := make([]string, 0)
	for _, dep := range app.Deps.Internal {
		if !checked[dep] && !app.Commands.Lint.ExcludeMap[dep] {
			flags = append(flags, "./"+dep)
			checking = append(checking, dep)
			log.Debug(dep)
		}
	}

	cmd := exec.Command("golangci-lint", flags...)
	cmd.Env = append(os.Environ(), app.Commands.Lint.EnvToList()...)

	if err := executeCommand(cmd); err != nil {
		return err
	}

	for _, dep := range checking {
		checked[dep] = true
	}
	return nil
}

func (app *App) Test(checked map[string]bool) error {
	flags := append([]string{"test"}, app.Commands.Test.Flags...)
	flags = append(flags, app.Location)

	checking := make([]string, 0)
	for _, dep := range app.Deps.Internal {
		if !checked[dep] && !app.Commands.Test.ExcludeMap[dep] {
			flags = append(flags, "./"+dep)
			checked[dep] = true
			log.Debug(dep)
		}
	}

	cmd := exec.Command("go", flags...)
	cmd.Env = append(os.Environ(), app.Commands.Test.EnvToList()...)

	if err := executeCommand(cmd); err != nil {
		return err
	}

	for _, dep := range checking {
		checked[dep] = true
	}
	return nil
}
