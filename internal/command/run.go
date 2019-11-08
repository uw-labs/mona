package command

import (
	"os"

	"github.com/apex/log"

	"github.com/uw-labs/mona/internal/app"
	"github.com/uw-labs/mona/internal/config"
	"github.com/uw-labs/mona/internal/execute"
)

// Run iterates over all new/modified apps and executes their lint, test and build commands
// as necessary. Once complete, appropriate hashes are updated in the lock file.
func Run(cfg Config) error {
	if err := Lint(cfg); err != nil {
		return err
	}
	if err := Test(cfg); err != nil {
		return err
	}
	return Build(cfg)
}

// Lint iterates over all new/modified apps and executes their lint command. Once complete,
// the lint hash is updated in the lock file.
func Lint(cfg Config) error {
	checked := make(map[string]bool)

	pkgs := cfg.Diff.PackagesList()
	log.Infof("Linting %v changed packages.", len(pkgs))
	if err := execute.Lint(config.CommandConfig{}, pkgs, checked); err != nil {
		return err
	}
	changed := getChangedAppNames(cfg)
	log.Infof("Linting %v changed apps.", len(changed))

	return execute.Lint(config.CommandConfig{}, changed, checked)
}

// Test attempts to run the test command for all apps where changes
// are detected.
func Test(cfg Config) error {
	checked := make(map[string]bool)

	pkgs := cfg.Diff.PackagesList()
	log.Infof("Testing %v changed packages.", len(pkgs))
	if err := execute.Test(config.CommandConfig{}, pkgs, checked); err != nil {
		return err
	}
	changed := getChangedAppNames(cfg)
	log.Infof("Testing %v changed apps.", len(changed))

	return execute.Test(config.CommandConfig{}, changed, checked)
}

// Build will execute the build commands for all apps where changes
// are detected.
func Build(cfg Config) error {
	if err := os.MkdirAll(cfg.Project.BinDir, os.ModePerm); err != nil {
		return err
	}

	return rangeChangedApps(cfg, func(app *app.App) error {
		log.Infof("Building app %s at %s", app.Name, app.Location)

		return app.Build("./" + cfg.Project.BinDir)
	})
}
