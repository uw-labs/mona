package files

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

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

// LoadProjectFile attempts to read a "mona.yml" file into memory from the provided
// working directory
func LoadProjectFile(wd string) (*ProjectFile, error) {
	file, err := os.Open(filepath.Join(wd, projectFileName))

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

// GetProjectRoot attempts to locate the root of the mona project based on the provided
// working directory. The provided path is traversed in reverse and each directory is
// checked for the existence of a mona.yml file.
func GetProjectRoot(wd string) (string, error) {
	sep := string(os.PathSeparator)
	parts := append(strings.SplitAfter(wd, sep), sep)

	for i := len(parts) - 1; i >= 0; i-- {
		dir := filepath.Join(parts[:i]...)
		file := filepath.Join(dir, projectFileName)

		if _, err := os.Stat(file); err == nil {
			return dir, nil
		} else if os.IsNotExist(err) {
			continue
		} else {
			return "", err
		}
	}

	return "", ErrNoProject
}
