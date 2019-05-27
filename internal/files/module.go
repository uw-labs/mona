package files

import (
	"errors"
	"os"
	"path"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

const (
	moduleFileName = "module.yml"
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
		Name      string   `yaml:"name"`                // The name of the module
		Exclude   []string `yaml:"exclude,omitempty"`   // File matchers to exclude files from hash generation.
		Artefacts []string `yaml:"artefacts,omitempty"` // Location of any artefacts that should be collected
		Location  string   `yaml:"-"`                   // The location of the module, not included in the module file but initialized externally for ease
		Commands  struct {
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

	if err != nil {
		return err
	}

	defer file.Close()
	mod := ModuleFile{
		Name: name,
	}

	return yaml.NewEncoder(file).Encode(mod)
}

// FindModules attempts to find all "module.yml" files in subdirectories of the given
// path and load them into memory.
func FindModules(dir string) (out []*ModuleFile, err error) {
	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
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

		out = append(out, module)
		return nil
	})

	return
}

// LoadModuleFile attempts to load a "module.yml" file into memory from
// the given location
func LoadModuleFile(location string) (*ModuleFile, error) {
	configLocation := filepath.Join(location, moduleFileName)
	file, err := os.Open(configLocation)

	if os.IsNotExist(err) {
		return nil, ErrNoModule
	}

	if err != nil {
		return nil, err
	}

	defer file.Close()

	var out ModuleFile

	if err := yaml.NewDecoder(file).Decode(&out); err != nil {
		return nil, err
	}

	out.Location = location
	return &out, nil
}

// UpdateModuleFile replaces the contents of "module.yml" at the given
// location with the module data provided.
func UpdateModuleFile(location string, module *ModuleFile) error {
	location = filepath.Join(location, moduleFileName)
	file, err := os.Create(location)

	if err != nil {
		return err
	}

	defer file.Close()
	return yaml.NewEncoder(file).Encode(module)
}

// CollectArtefacts attempts to move all files specified in the module to the provided
// output directory. Files are stored in a folder named after the module.
func (m *ModuleFile) CollectArtefacts(outputDir string) error {
	if err := os.MkdirAll(path.Join(outputDir, m.Name), 0777); err != nil {
		return err
	}

	for _, artefactName := range m.Artefacts {
		dest := path.Join(outputDir, m.Name, artefactName)
		src := path.Join(m.Location, artefactName)

		if err := os.Rename(src, dest); err != nil {
			return err
		}
	}

	return nil
}
