package execute

import (
	"os"
	"os/exec"

	"github.com/uw-labs/mona/internal/config"
)

func Lint(cfg config.CommandConfig, pkgs []string, checked map[string]bool) error {
	flags := append([]string{"run"}, cfg.Flags...)

	checking := make([]string, 0, len(pkgs))

	for _, pkg := range pkgs {
		if !checked[pkg] {
			flags = append(flags, "./"+pkg)
			checking = append(checking, pkg)
		}
	}

	cmd := exec.Command("golangci-lint", flags...)
	cmd.Env = append(os.Environ(), cfg.EnvToList()...)

	if err := Command(cmd); err != nil {
		return err
	}

	for _, pkg := range checking {
		checked[pkg] = true
	}
	return nil
}
