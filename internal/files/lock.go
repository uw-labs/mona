package files

import (
	"errors"
	"os"

	"gopkg.in/yaml.v2"
)

const (
	lockFileName = "mona.lock"
)

var (
	// ErrNoLock is the error returned when no lock file is present in the current
	// working directory.
	ErrNoLock = errors.New(`no "mona.lock" file in current wd`)
)

type (
	// The Lock type represents the structure of a lock file, it stores the project name,
	// version and the last build hashes used for each module
	Lock struct {
		Name    string          `yaml:"name"`
		Version string          `yaml:"version"`
		Modules []ModuleVersion `yaml:"modules,omitempty"`
	}

	// The ModuleVersion type represents individual module information as stored
	// in the lock file.
	ModuleVersion struct {
		Name      string `yaml:"name"`
		Location  string `yaml:"location"`
		BuildHash string `yaml:"build"`
	}
)

// NewLockFile creates a new "mona.lock" file at the project root using the provided name
// and version
func NewLockFile(name, version string) error {
	file, err := os.Create(lockFileName)

	if err != nil {
		return err
	}

	defer file.Close()
	lock := Lock{
		Name:    name,
		Version: version,
	}

	return yaml.NewEncoder(file).Encode(lock)
}

// UpdateLockFile overwrites the current "mona.lock" file at the project root with
// the data provided.
func UpdateLockFile(lock *Lock) error {
	file, err := os.Create(lockFileName)

	if err != nil {
		return err
	}

	defer file.Close()
	return yaml.NewEncoder(file).Encode(lock)
}

// LoadLockFile reads the "mona.lock" file from the project root into
// memory.
func LoadLockFile() (*Lock, error) {
	file, err := os.Open(lockFileName)

	if os.IsNotExist(err) {
		return nil, ErrNoLock
	}

	if err != nil {
		return nil, err
	}

	defer file.Close()

	var out Lock

	if err := yaml.NewDecoder(file).Decode(&out); err != nil {
		return nil, err
	}

	return &out, nil
}

// AddModule adds a new module entry to the lock file using the provided name,
// location and hash. The "mona.lock" file is then updated with the new values.
func (l *Lock) AddModule(name, location, hash string) error {
	l.Modules = append(l.Modules, ModuleVersion{
		Name:      name,
		Location:  location,
		BuildHash: hash,
	})

	file, err := os.Create(lockFileName)

	if err != nil {
		return err
	}

	defer file.Close()
	return yaml.NewEncoder(file).Encode(l)
}
