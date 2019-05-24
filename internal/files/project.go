package files

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

const (
	projectFileName = "mona.yml"
)

type (
	// The ProjectFile type represents the structure of the "mona.yml" file.
	ProjectFile struct {
		Name    string   `yaml:"name"`              // The name of the project
		Version string   `yaml:"version"`           // The mona version used to create the project
		Modules []string `yaml:"modules,omitempty"` // The modules used within the project.
	}
)

// NewProjectFile creates a new "mona.yml" file in the current working directory with the given
// name and version.
func NewProjectFile(name, version string) error {
	file, err := os.Create(projectFileName)

	if err != nil {
		return err
	}

	defer file.Close()
	pj := ProjectFile{
		Name:    name,
		Version: version,
	}

	return yaml.NewEncoder(file).Encode(pj)
}

// LoadProjectFile attempts to read a "mona.yml" file into memory.
func LoadProjectFile() (*ProjectFile, error) {
	file, err := os.Open(projectFileName)

	if err != nil {
		return nil, err
	}

	defer file.Close()

	var out ProjectFile

	if err := yaml.NewDecoder(file).Decode(&out); err != nil {
		return nil, err
	}

	return &out, nil
}

// AddModule writes a new module to the "mona.yml" file.
func (p *ProjectFile) AddModule(name, location string) error {
	for _, mod := range p.Modules {
		parts := strings.Split(mod, " ")

		if parts[0] == name && parts[1] == location {
			return fmt.Errorf("module %s in %s already exists", name, location)
		}
	}

	p.Modules = append(p.Modules, fmt.Sprintf("%s %s", name, location))

	file, err := os.Create(projectFileName)

	if err != nil {
		return err
	}

	defer file.Close()
	return yaml.NewEncoder(file).Encode(p)
}
