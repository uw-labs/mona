package files

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/hashicorp/go-multierror"
	"github.com/iafan/cwalk"
	"gopkg.in/yaml.v2"
)

const (
	moduleFileName = "module.yml"
	moduleFilePerm = 0777
)

var (
	// ErrNoModule is the error returned when a module for a requested directory
	// does not exist.
	ErrNoModule = errors.New("could not find module.yml at the specified location")
)

type (
	// The ModuleFile type represents the data held in the "module.yml" file in each module
	// directory
	ModuleFile struct {
		Name     string   `yaml:"name"`              // The name of the module
		Exclude  []string `yaml:"exclude,omitempty"` // File matchers to exclude files from hash generation.
		Location string   `yaml:"-"`                 // The location of the module, not included in the module file but initialized externally for ease
		Commands struct {
			Build string `yaml:"build"` // Command for building the module
			Test  string `yaml:"test"`  // Command for testing the module
			Lint  string `yaml:"lint"`  // Command for linting the module
		} `yaml:"commands"` // Commands that can be executed against the module
	}
)

// NewModuleFile creates a new "module.yml" file with the given name at the given
// location.
func NewModuleFile(name, location string) error {
	location = filepath.Join(location, moduleFileName)
	file, err := os.Create(location)

	if os.IsNotExist(err) {
		return fmt.Errorf("directory %s does not exist", location)
	}

	if err != nil {
		return err
	}

	mod := ModuleFile{
		Name: name,
	}

	return multierror.Append(
		yaml.NewEncoder(file).Encode(mod),
		file.Close()).
		ErrorOrNil()
}

// FindModules attempts to find all "module.yml" files in subdirectories of the given
// path and load them into memory.
func FindModules(dir string) (out []*ModuleFile, err error) {
	var mux sync.Mutex
	err = cwalk.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if info.Name() != moduleFileName {
			return nil
		}

		module, err := LoadModuleFile(strings.TrimSuffix(path, moduleFileName))

		if err != nil {
			return err
		}

		mux.Lock()
		defer mux.Unlock()
		out = append(out, module)

		return filepath.SkipDir
	})

	return
}

// LoadModuleFile attempts to load a "module.yml" file into memory from
// the given location
func LoadModuleFile(location string) (*ModuleFile, error) {
	file, err := os.OpenFile(
		filepath.Join(location, moduleFileName),
		os.O_RDONLY,
		moduleFilePerm)

	if os.IsNotExist(err) {
		return nil, ErrNoModule
	}

	if err != nil {
		return nil, err
	}

	var out ModuleFile

	if err := yaml.NewDecoder(file).Decode(&out); err != nil {
		return nil, err
	}

	out.Location = location
	return &out, file.Close()
}

// UpdateModuleFile replaces the contents of "module.yml" at the given
// location with the module data provided.
func UpdateModuleFile(location string, module *ModuleFile) error {
	file, err := os.OpenFile(
		filepath.Join(location, moduleFileName),
		os.O_WRONLY,
		moduleFilePerm)

	if os.IsNotExist(err) {
		return ErrNoModule
	}

	if err != nil {
		return err
	}

	return multierror.Append(
		yaml.NewEncoder(file).Encode(module),
		file.Close()).
		ErrorOrNil()
}
