// Package command contains methods called by the CLI to manage
// a mona project.
package command

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
	"strings"

	"github.com/apex/log"
	"github.com/davidsbond/mona/internal/config"
	"github.com/davidsbond/mona/internal/deps"
	"github.com/davidsbond/mona/internal/hash"
	"github.com/hashicorp/go-multierror"
)

type (
	changeType int
	rangeFn    func(*config.AppFile) error
)

const (
	changeTypeBuild changeType = 0
	changeTypeTest  changeType = 1
	changeTypeLint  changeType = 2
)

func getChangedApps(mod deps.Module, pj *config.ProjectFile, change changeType) ([]*config.AppFile, error) {
	apps, err := config.FindApps(pj.Location)

	if err != nil {
		return nil, err
	}

	lock, err := config.LoadLockFile(pj.Location)

	if err != nil {
		return nil, err
	}

	var out []*config.AppFile
	for _, appInfo := range apps {
		lockInfo, ok := lock.Apps[appInfo.Name]

		if !ok {
			// If the app isn't present in the lock file, it's probably
			// new and needs adding to the lock file.
			out = append(out, appInfo)
			continue
		}

		// GenerateString a new hash for the app directory
		exclude := append(pj.Exclude, appInfo.Exclude...)
		newHash, err := hash.GetForApp(mod, appInfo.Location, exclude...)

		if err != nil {
			return nil, err
		}

		// Check if we need to build/lint/test
		var oldHash string
		switch change {
		case changeTypeBuild:
			oldHash = lockInfo.BuildHash
		case changeTypeTest:
			oldHash = lockInfo.TestHash
		case changeTypeLint:
			oldHash = lockInfo.LintHash
		}

		if oldHash != newHash {
			out = append(out, appInfo)
		}
	}

	return out, nil
}

func streamOutputs(outputs ...io.ReadCloser) {
	for _, output := range outputs {
		go func(o io.ReadCloser) {
			defer o.Close()

			scanner := bufio.NewScanner(o)
			scanner.Split(bufio.ScanLines)

			for scanner.Scan() {
				m := scanner.Text()
				fmt.Println(m)
			}
		}(output)
	}
}

func streamCommand(command, wd string) error {
	parts := strings.Split(command, " ")
	cmd := exec.Command(parts[0], parts[1:]...)
	cmd.Dir = wd

	stdout, err := cmd.StdoutPipe()

	if err != nil {
		return err
	}

	stderr, err := cmd.StderrPipe()

	if err != nil {
		return err
	}

	streamOutputs(stdout, stderr)

	if err := cmd.Start(); err != nil {
		return err
	}

	return cmd.Wait()
}

func rangeChangedApps(mod deps.Module, pj *config.ProjectFile, fn rangeFn, ct changeType) error {
	changed, err := getChangedApps(mod, pj, ct)

	if err != nil || len(changed) == 0 {
		return err
	}

	lock, err := config.LoadLockFile(pj.Location)

	if err != nil {
		return err
	}

	var errs []error
	for _, app := range changed {
		if err := fn(app); err != nil {
			errs = append(errs, fmt.Errorf("app %s: %s", app.Name, err.Error()))
			continue
		}

		exclude := append(pj.Exclude, app.Exclude...)
		newHash, err := hash.GetForApp(mod, app.Location, exclude...)

		if err != nil {
			return err
		}

		lockInfo, modInLock := lock.Apps[app.Name]

		if !modInLock {
			log.Debugf("Detected new app %s at %s, adding to lock file", app.Name, app.Location)

			if err := config.AddApp(lock, pj.Location, app.Name); err != nil {
				return err
			}

			lockInfo = lock.Apps[app.Name]
		}

		switch ct {
		case changeTypeBuild:
			lockInfo.BuildHash = newHash
		case changeTypeTest:
			lockInfo.TestHash = newHash
		case changeTypeLint:
			lockInfo.LintHash = newHash
		}

		if err := config.UpdateLockFile(pj.Location, lock); err != nil {
			return err
		}
	}

	return multierror.Append(nil, errs...).ErrorOrNil()
}
