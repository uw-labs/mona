package files

import (
	"os"
	"path"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

const (
	moduleFileName = "module.yml"
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
			Build string `yaml:"build"` // Commands for building the module
			Test  string `yaml:"test"`  // Commands for testing the module
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

// LoadModuleFile attempts to load a "module.yml" file into memory from
// the given location
func LoadModuleFile(location string) (*ModuleFile, error) {
	configLocation := filepath.Join(location, moduleFileName)
	file, err := os.Open(configLocation)

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
	if err := os.MkdirAll(outputDir, 0777); err != nil {
		return err
	}

	for _, artefactName := range m.Artefacts {
		if err := os.MkdirAll(path.Join(outputDir, m.Name), 0777); err != nil {
			return err
		}

		dest := path.Join(outputDir, m.Name, artefactName)
		src := path.Join(m.Location, artefactName)

		if err := os.Rename(src, dest); err != nil {
			return err
		}
	}

	return nil
}
