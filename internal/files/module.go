package files

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

const (
	moduleFileName = "module.yml"
)

type (
	// The Module type represents the data held in the "module.yml" file in each module
	// directory
	Module struct {
		Name     string `yaml:"name"` // The name of the module
		Commands struct {
			Build string `yaml:"build"` // Commands for building the module
			Test  string `yaml:"test"`  // Commands for testing the module
		} `yaml:"commands"` // Commands that can be executed against the module
		Exclude  []string `yaml:"exclude,omitempty"`
		Location string   `yaml:"-"` // The location of the module, not included in the module file but initialized externally for ease
	}
)

// NewModuleFile creates a new "module.yml" file with the given name at the given
// location.
func NewModuleFile(name, location string) error {
	location = filepath.Join(location, moduleFileName)
	file, err := os.Create(location)

	if err != nil {
		return err
	}

	defer file.Close()
	mod := Module{
		Name: name,
	}

	return yaml.NewEncoder(file).Encode(mod)
}

// LoadModuleFile attempts to load a "module.yml" file into memory from
// the given location
func LoadModuleFile(location string) (*Module, error) {
	location = filepath.Join(location, moduleFileName)
	file, err := os.Open(location)

	if err != nil {
		return nil, err
	}

	defer file.Close()

	var out Module

	if err := yaml.NewDecoder(file).Decode(&out); err != nil {
		return nil, err
	}

	return &out, nil
}
