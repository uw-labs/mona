// Package command contains methods called by the CLI to manage
// a mona project.
package command

import (
	"fmt"

	"github.com/apex/log"
	"github.com/hashicorp/go-multierror"

	"github.com/davidsbond/mona/internal/app"
	"github.com/davidsbond/mona/internal/config"
	"github.com/davidsbond/mona/internal/hash"
)

type Config struct {
	Project  *config.Project
	FailFast bool
}

type (
	changeType int
	rangeFn    func(*app.App, map[changeType]bool) error
	changedApp struct {
		app         *app.App
		newHash     string
		changeTypes map[changeType]bool
	}
)

const (
	changeTypeLint changeType = iota
	changeTypeTest
	changeTypeBuild
)

func getLockAndChangedApps(pj *config.Project) (lock *config.LockFile, out []changedApp, err error) {
	apps, err := app.FindApps("./", pj.Mod)
	if err != nil {
		return nil, nil, err
	}

	lock, err = config.LoadLockFile(pj.Location)
	if err != nil {
		return nil, nil, err
	}

	for _, appInfo := range apps {
		lockInfo := lock.Apps[appInfo.Name]

		// GenerateString a new hash for the app directory
		exclude := append(pj.Exclude, appInfo.Exclude...)
		newHash, err := hash.GetAppDeps(pj.Mod, appInfo.Location, exclude...)
		if err != nil {
			return nil, nil, err
		}

		changes := make(map[changeType]bool, 3)

		if lockInfo.LintHash != newHash {
			changes[changeTypeLint] = true
		}
		if lockInfo.TestHash != newHash {
			changes[changeTypeTest] = true
		}
		if lockInfo.BuildHash != newHash {
			changes[changeTypeBuild] = true
		}
		if len(changes) != 0 {
			out = append(out, changedApp{
				app:         appInfo,
				newHash:     newHash,
				changeTypes: changes,
			})
		}
	}

	return lock, out, nil
}

func rangeChangedApps(cfg Config, cts []changeType, fn rangeFn) error {
	lock, changed, err := getLockAndChangedApps(cfg.Project)

	if err != nil || len(changed) == 0 {
		return err
	}

	var errs []error
	for _, change := range changed {
		if len(cts) == 1 && !change.changeTypes[cts[0]] {
			continue
		}
		if err := fn(change.app, change.changeTypes); err != nil {
			errs = append(errs, fmt.Errorf("app %s: %s", change.app.Name, err.Error()))

			if cfg.FailFast {
				return errs[0]
			}
			continue
		}

		exclude := append(cfg.Project.Exclude, change.app.Exclude...)
		newHash, err := hash.GetAppDeps(cfg.Project.Mod, change.app.Location, exclude...)

		if err != nil {
			return err
		}

		lockInfo, modInLock := lock.Apps[change.app.Name]

		if !modInLock {
			log.Debugf("Detected new appInfo %s at %s, adding to lock file", change.app.Name, change.app.Location)

			if err := config.AddApp(lock, cfg.Project.Location, change.app.Name); err != nil {
				return err
			}

			lockInfo = lock.Apps[change.app.Name]
		}

		for _, ct := range cts {
			switch ct {
			case changeTypeLint:
				lockInfo.LintHash = newHash
			case changeTypeTest:
				lockInfo.TestHash = newHash
			case changeTypeBuild:
				lockInfo.BuildHash = newHash
			}
		}

		if err := config.UpdateLockFile(cfg.Project.Location, lock); err != nil {
			return err
		}
	}

	return multierror.Append(nil, errs...).ErrorOrNil()
}
