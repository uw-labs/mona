package config

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/apex/log"
	"github.com/hashicorp/go-multierror"
	"gopkg.in/yaml.v2"

	"github.com/uw-labs/mona/internal/golang"
)

const (
	goModFileName   = "go.mod"
	projectFileName = "mona.yml"
	projectFilePerm = 0777
)

var (
	// ErrNoProject is the error returned when a project file is not
	// found in the current directory
	ErrNoProject = errors.New("failed to find project file in current path")
	// ErrNoGoModFile is an error signifying that a directory where
	// the user wants to initialise a project doesn't contain go.mod file.
	ErrNoGoModFile = errors.New("no go.mod file, run 'go mod init' first")
)

type (
	// Project type represents the structure of the "mona.yml" file.
	Project struct {
		Name     string        `yaml:"name"`              // The name of the project
		Exclude  []string      `yaml:"exclude,omitempty"` // Global file patterns to ignore during hash generation
		Location string        `yaml:"-"`                 // The root project directory, not set in the yaml file but set on load for convenience
		BinDir   string        `yaml:"binDir"`            // Relative path from the root of the project to the directory where compiled binaries will be placed.
		Mod      golang.Module `yaml:"-"`
		Branch   string        `yaml:"-"`
	}
)

// NewProject creates a new "mona.yml" file in the provided directory with the given name.
func NewProject(dir string, name string) error {
	if _, err := os.Stat(filepath.Join(dir, goModFileName)); err != nil {
		if os.IsNotExist(err) {
			return ErrNoGoModFile
		}
		return err
	}
	location := filepath.Join(dir, projectFileName)
	file, err := os.Create(location)

	if err != nil {
		return err
	}

	pj := Project{
		Name:   name,
		BinDir: "bin",
	}

	return multierror.Append(
		yaml.NewEncoder(file).Encode(pj),
		file.Close()).
		ErrorOrNil()
}

// LoadProject attempts to read a "mona.yml" file into memory from the provided
// working directory
func LoadProject(wd, branch string) (*Project, error) {
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
	defer file.Close()

	var out Project
	if err := yaml.NewDecoder(file).Decode(&out); err != nil {
		return nil, err
	}

	out.Mod, err = golang.ParseModuleFile(filepath.Join(wd, "go.mod"))
	if err != nil {
		return nil, err
	}
	out.Location = wd
	out.Branch = branch

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
			log.Debugf("Found project file at %s", dir)
			return dir, nil
		} else if os.IsNotExist(err) {
			continue
		} else {
			return "", err
		}
	}

	return "", ErrNoProject
}
