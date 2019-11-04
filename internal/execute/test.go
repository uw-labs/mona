package execute

import (
	"os"
	"os/exec"

	"github.com/uw-labs/mona/internal/config"
)

func Test(cfg config.CommandConfig, pkgs []string, checked map[string]bool) error {
	flags := append([]string{"test"}, cfg.Flags...)
	checking := make([]string, 0, len(pkgs))

	for _, pkg := range pkgs {
		if !checked[pkg] {
			flags = append(flags, "./"+pkg)
			checking = append(checking, pkg)
		}
	}
	cmd := exec.Command("go", flags...)
	cmd.Env = append(os.Environ(), cfg.EnvToList()...)

	if err := Command(cmd); err != nil {
		return err
	}

	for _, pkg := range checking {
		checked[pkg] = true
	}
	return nil
}
