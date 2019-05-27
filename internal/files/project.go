package files

import (
	"errors"
	"os"

	"gopkg.in/yaml.v2"
)

const (
	projectFileName = "mona.yml"
)

var (
	// ErrNoProject is the error returned when a project file is not
	// found in the current directory
	ErrNoProject = errors.New("no mona.yml file found in directory")
)

type (
	// The ProjectFile type represents the structure of the "mona.yml" file.
	ProjectFile struct {
		Name      string `yaml:"name"`                // The name of the project
		Version   string `yaml:"version"`             // The mona version used to create the project
		Artefacts string `yaml:"artefacts,omitempty"` // The location for artefacts to be stored
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

	if os.IsNotExist(err) {
		return nil, ErrNoProject
	}

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
