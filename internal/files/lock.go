package files

import (
	"errors"
	"os"
	"path/filepath"

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
	// The LockFile type represents the structure of a lock file, it stores the project name,
	// version and the last build hashes used for each module
	LockFile struct {
		Name    string                    `yaml:"name"`
		Version string                    `yaml:"version"`
		Modules map[string]*ModuleVersion `yaml:"modules,omitempty"`
	}

	// The ModuleVersion type represents individual module information as stored
	// in the lock file.
	ModuleVersion struct {
		BuildHash string `yaml:"build"`
		TestHash  string `yaml:"test"`
		LintHash  string `yaml:"lint"`
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
	lock := LockFile{
		Name:    name,
		Version: version,
		Modules: make(map[string]*ModuleVersion),
	}

	return yaml.NewEncoder(file).Encode(lock)
}

// UpdateLockFile overwrites the current "mona.lock" file at the project root with
// the data provided.
func UpdateLockFile(lock *LockFile) error {
	file, err := os.Create(lockFileName)

	if err != nil {
		return err
	}

	defer file.Close()
	return yaml.NewEncoder(file).Encode(lock)
}

// LoadLockFile attempts to load a lock file into memory from the provided
// working directory.
func LoadLockFile(wd string) (*LockFile, error) {
	file, err := os.Open(filepath.Join(wd, lockFileName))

	if os.IsNotExist(err) {
		return nil, ErrNoLock
	}

	if err != nil {
		return nil, err
	}

	defer file.Close()

	var out LockFile

	if err := yaml.NewDecoder(file).Decode(&out); err != nil {
		return nil, err
	}

	if out.Modules == nil {
		out.Modules = make(map[string]*ModuleVersion)
	}

	return &out, nil
}

// AddModule adds a new module entry to the lock file using the provided name,
// location and hash. The "mona.lock" file is then updated with the new values.
func (l *LockFile) AddModule(name string) error {
	l.Modules[name] = &ModuleVersion{}

	file, err := os.Create(lockFileName)

	if err != nil {
		return err
	}

	defer file.Close()
	return yaml.NewEncoder(file).Encode(l)
}
