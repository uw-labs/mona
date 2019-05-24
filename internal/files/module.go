package files

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

const (
	moduleFileName = "module.yaml"
)

type (
	Module struct {
		Name     string `yaml:"name"`
		Commands struct {
			Build string `yaml:"build"`
			Test  string `yaml:"test"`
		} `yaml:"commands"`
		Exclude  []string `yaml:"exclude,omitempty"`
		Location string   `yaml:"-"` // The location of the module, not included in the module file but initialized externally for ease
	}
)

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
