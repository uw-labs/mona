package files

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

const (
	projectFileName = "mona.yml"
)

type (
	Project struct {
		Name    string   `yaml:"name"`
		Version string   `yaml:"version"`
		Modules []string `yaml:"modules,omitempty"`
	}
)

func NewProjectFile(name, version string) error {
	file, err := os.Create(projectFileName)

	if err != nil {
		return err
	}

	defer file.Close()
	pj := Project{
		Name:    name,
		Version: version,
	}

	return yaml.NewEncoder(file).Encode(pj)
}

func LoadProjectFile() (*Project, error) {
	file, err := os.Open(projectFileName)

	if err != nil {
		return nil, err
	}

	defer file.Close()

	var out Project

	if err := yaml.NewDecoder(file).Decode(&out); err != nil {
		return nil, err
	}

	return &out, nil
}

func (p *Project) AddModule(name, location string) error {
	p.Modules = append(p.Modules, fmt.Sprintf("%s %s", name, location))

	file, err := os.Create(projectFileName)

	if err != nil {
		return err
	}

	defer file.Close()
	return yaml.NewEncoder(file).Encode(p)
}
