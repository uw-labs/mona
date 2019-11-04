package app

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/apex/log"
	"github.com/hashicorp/go-multierror"
	"github.com/iafan/cwalk"
	"gopkg.in/yaml.v2"

	"github.com/davidsbond/mona/internal/config"
	"github.com/davidsbond/mona/internal/deps"
)

const (
	appFileName = "app.yml"
	appFilePerm = 0777
)

var (
	// ErrNoApp is the error returned when an app for a requested directory
	// does not exist.
	ErrNoApp = errors.New("could not find app.yml at the specified location")
)

type (
	// The App a configuration for a single app that is held in the "app.yml" file.
	App struct {
		Name     string   `yaml:"name"`              // The name of the app
		Exclude  []string `yaml:"exclude,omitempty"` // File matchers to exclude files from hash generation.
		Location string   `yaml:"-"`                 // The location of the app, not included in the app file but initialized externally for ease
		Commands struct {
			Build config.BuildConfig   // Command for building the app
			Test  config.CommandConfig // Command for testing the app
			Lint  config.CommandConfig // Command for linting the app
		} `yaml:"commands"` // Commands that can be executed against the app
		Deps deps.AppDeps `yaml:"-"`
	}
)

// NewAppFile creates a new "app.yml" file with the given name at the given location.
func NewAppFile(name, location string) error {
	location = filepath.Join(location, appFileName)
	file, err := os.Create(location)

	if os.IsNotExist(err) {
		return fmt.Errorf("directory %s does not exist", location)
	}

	if err != nil {
		return err
	}

	mod := App{
		Name: name,
	}

	return multierror.Append(
		yaml.NewEncoder(file).Encode(mod),
		file.Close(),
	).ErrorOrNil()
}

// FindApps attempts to find all "app.yml" files in subdirectories of the given
// path and load them into memory.
func FindApps(projectDir string, mod deps.Module) (out []*App, outErr error) {
	dir := projectDir
	log.Debugf("Searching for apps in %s", dir)

	var (
		appMux  sync.Mutex
		errMux  sync.Mutex
		skipMux sync.Mutex
		skip    []string
	)

	walkErr := cwalk.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			errMux.Lock()
			outErr = multierror.Append(outErr, err)
			errMux.Unlock()
			return err
		}

		skipMux.Lock()
		for _, s := range skip {
			if strings.HasPrefix(path, s) {
				skipMux.Unlock()
				return filepath.SkipDir
			}
		}
		skipMux.Unlock()

		if info.IsDir() {
			return nil
		}

		if info.Name() != appFileName {
			return nil
		}

		app, err := LoadApp(strings.TrimSuffix(path, appFileName), mod)
		if err != nil {
			errMux.Lock()
			outErr = multierror.Append(outErr, err)
			errMux.Unlock()
			return err
		}

		log.Debugf("Found app %s at %s", app.Name, app.Location)

		appMux.Lock()
		out = append(out, app)
		appMux.Unlock()

		dir, _ := filepath.Split(path)

		skipMux.Lock()
		skip = append(skip, dir)
		skipMux.Unlock()

		return filepath.SkipDir
	})
	if walkErr != nil {
		outErr = multierror.Append(outErr, walkErr)
	}

	return out, outErr
}

// LoadApp attempts to load a "app.yml" file into memory from
// the given location.
func LoadApp(location string, mod deps.Module) (*App, error) {
	file, err := os.OpenFile(
		filepath.Join(location, appFileName),
		os.O_RDONLY,
		appFilePerm)

	if os.IsNotExist(err) {
		return nil, ErrNoApp
	}
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var out App
	if err := yaml.NewDecoder(file).Decode(&out); err != nil {
		return nil, err
	}

	out.Commands.Lint.ExcludeMap = make(map[string]bool, len(out.Commands.Lint.Exclude))
	for _, exclude := range out.Commands.Lint.Exclude {
		out.Commands.Lint.ExcludeMap[exclude] = true
	}
	out.Commands.Test.ExcludeMap = make(map[string]bool, len(out.Commands.Test.Exclude))
	for _, exclude := range out.Commands.Test.Exclude {
		out.Commands.Test.ExcludeMap[exclude] = true
	}

	out.Location = "./" + location
	out.Deps, err = deps.GetAppDeps(mod, out.Location)
	if err != nil {
		return nil, err
	}

	return &out, nil
}

// SaveApp replaces the contents of "app.yml" at the given
// location with the app data provided.
func SaveApp(location string, app *App) error {
	file, err := os.OpenFile(
		filepath.Join(location, appFileName),
		os.O_WRONLY,
		appFilePerm)

	if os.IsNotExist(err) {
		return ErrNoApp
	}

	if err != nil {
		return err
	}

	return multierror.Append(
		yaml.NewEncoder(file).Encode(app),
		file.Close()).
		ErrorOrNil()
}
