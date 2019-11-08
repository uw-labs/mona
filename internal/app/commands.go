package app

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/apex/log"

	"github.com/uw-labs/mona/internal/execute"
)

func (app *App) Build(binDir string) error {
	flags := append([]string{"build"}, fmt.Sprintf("-o=%s/%s", binDir, app.Name))
	flags = append(flags, app.Commands.Build.AllFlags()...)
	flags = append(flags, "./"+app.Location)

	cmd := exec.Command("go", flags...)
	cmd.Env = append(os.Environ(), app.Commands.Build.EnvToList()...)
	log.Info(cmd.String())

	return execute.Command(cmd)
}
