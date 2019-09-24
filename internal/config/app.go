package config

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
	// The AppFile type represents the data held in the "app.yml" file in each app
	// directory
	AppFile struct {
		Name     string   `yaml:"name"`              // The name of the app
		Exclude  []string `yaml:"exclude,omitempty"` // File matchers to exclude files from hash generation.
		Location string   `yaml:"-"`                 // The location of the app, not included in the app file but initialized externally for ease
		Commands struct {
			Build struct {
				Run   string `yaml:"run"`   // The command to run
				Image string `yaml:"image"` // The docker image to use
			} `yaml:"build"` // Command for building the app
			Test struct {
				Run   string `yaml:"run"`   // The command to run
				Image string `yaml:"image"` // The docker image to use
			} `yaml:"test"` // Command for testing the app
			Lint struct {
				Run   string `yaml:"run"`   // The command to run
				Image string `yaml:"image"` // The docker image to use
			} `yaml:"lint"` // Command for linting the app
		} `yaml:"commands"` // Commands that can be executed against the app
	}
)

// NewAppFile creates a new "app.yml" file with the given name at the given
// location.
func NewAppFile(name, location string) error {
	location = filepath.Join(location, appFileName)
	file, err := os.Create(location)

	if os.IsNotExist(err) {
		return fmt.Errorf("directory %s does not exist", location)
	}

	if err != nil {
		return err
	}

	mod := AppFile{
		Name: name,
	}

	return multierror.Append(
		yaml.NewEncoder(file).Encode(mod),
		file.Close()).
		ErrorOrNil()
}

// FindApps attempts to find all "app.yml" files in subdirectories of the given
// path and load them into memory.
func FindApps(dir string) (out []*AppFile, err error) {
	log.Debugf("Searching for apps in %s", dir)

	var appMux sync.Mutex
	var skipMux sync.Mutex
	var skip []string

	err = cwalk.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
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

		app, err := LoadAppFile(strings.TrimSuffix(path, appFileName))

		if err != nil {
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

	return
}

// LoadAppFile attempts to load a "app.yml" file into memory from
// the given location
func LoadAppFile(location string) (*AppFile, error) {
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

	var out AppFile

	if err := yaml.NewDecoder(file).Decode(&out); err != nil {
		return nil, err
	}

	out.Location = location
	return &out, file.Close()
}

// UpdateAppFile replaces the contents of "app.yml" at the given
// location with the app data provided.
func UpdateAppFile(location string, app *AppFile) error {
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
