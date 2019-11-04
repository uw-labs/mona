// Package command contains methods called by the CLI to manage
// a mona project.
package command

import (
	"fmt"

	"github.com/hashicorp/go-multierror"

	"github.com/uw-labs/mona/internal/app"
)

type (
	changeType int
	rangeFn    func(*app.App) error
)

const (
	changeTypeLint changeType = iota
	changeTypeTest
	changeTypeBuild
)

func getChangedApps(cfg Config) []*app.App {
	out := make([]*app.App, 0, len(cfg.Apps))

LOOP:
	for _, appInfo := range cfg.Apps {
		if cfg.Diff.Packages[appInfo.Location] {
			out = append(out, appInfo)
			continue
		}
		for _, pkg := range appInfo.Deps.Internal {
			if cfg.Diff.Packages[pkg] {
				out = append(out, appInfo)
				continue LOOP
			}
		}
		for _, mod := range appInfo.Deps.External {
			if cfg.Diff.Modules[mod] {
				out = append(out, appInfo)
				continue LOOP
			}
		}
	}

	return out
}

func getChangedAppNames(cfg Config) (out []string) {
	for _, appInfo := range getChangedApps(cfg) {
		out = append(out, appInfo.Location)
	}
	return out
}

func rangeChangedApps(cfg Config, fn rangeFn) error {
	var errs []error
	for _, appInfo := range getChangedApps(cfg) {
		if err := fn(appInfo); err != nil {
			errs = append(errs, fmt.Errorf("app %s: %s", appInfo.Name, err.Error()))

			if cfg.FailFast {
				return errs[0]
			}
			continue
		}
	}

	return multierror.Append(nil, errs...).ErrorOrNil()
}
