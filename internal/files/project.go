package files

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/hashicorp/go-multierror"

	"gopkg.in/yaml.v2"
)

const (
	projectFileName = "mona.yml"
	projectFilePerm = 0777
)

var (
	// ErrNoProject is the error returned when a project file is not
	// found in the current directory
	ErrNoProject = errors.New("failed to find project file in current path")
)

type (
	// The ProjectFile type represents the structure of the "mona.yml" file.
	ProjectFile struct {
		Name     string   `yaml:"name"`              // The name of the project
		Exclude  []string `yaml:"exclude,omitempty"` // Global file patterns to ignore during hash generation
		Location string   `yaml:"-"`                 // The root project directory, not set in the yaml file but set on load for convenience
	}
)

// NewProjectFile creates a new "mona.yml" file in the current working directory with the given
// name.
func NewProjectFile(name string) error {
	file, err := os.Create(projectFileName)

	if err != nil {
		return err
	}

	pj := ProjectFile{
		Name: name,
	}

	return multierror.Append(
		yaml.NewEncoder(file).Encode(pj),
		file.Close()).
		ErrorOrNil()
}

// LoadProjectFile attempts to read a "mona.yml" file into memory from the provided
// working directory
func LoadProjectFile(wd string) (*ProjectFile, error) {
	file, err := os.OpenFile(
		filepath.Join(wd, projectFileName),
		os.O_RDONLY,
		projectFilePerm)

	if os.IsNotExist(err) {
		return nil, ErrNoProject
	}

	if err != nil {
		return nil, err
	}

	var out ProjectFile

	if err := yaml.NewDecoder(file).Decode(&out); err != nil {
		return nil, err
	}

	out.Location = wd
	return &out, file.Close()
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
