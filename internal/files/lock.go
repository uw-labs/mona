package files

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/hashicorp/go-multierror"
	"gopkg.in/yaml.v2"
)

const (
	lockFileName = "mona.lock"
	lockFilePerm = 0644
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

// NewLockFile creates a new "mona.lock" file in the current working directory using the
// provided name and version
func NewLockFile(name, version string) error {
	file, err := os.Create(lockFileName)

	if err != nil {
		return err
	}

	lock := LockFile{
		Name:    name,
		Version: version,
		Modules: make(map[string]*ModuleVersion),
	}

	return multierror.Append(
		yaml.NewEncoder(file).Encode(lock),
		file.Close()).
		ErrorOrNil()
}

// UpdateLockFile overwrites the current "mona.lock" file in the given working
// directory with the data provided.
func UpdateLockFile(wd string, lock *LockFile) error {
	file, err := os.OpenFile(
		filepath.Join(wd, lockFileName),
		os.O_WRONLY,
		lockFilePerm)

	if os.IsNotExist(err) {
		return ErrNoLock
	}

	if err != nil {
		return err
	}

	return multierror.Append(
		yaml.NewEncoder(file).Encode(lock),
		file.Close()).
		ErrorOrNil()
}

// LoadLockFile attempts to load a lock file into memory from the provided
// working directory.
func LoadLockFile(wd string) (*LockFile, error) {
	file, err := os.OpenFile(
		filepath.Join(wd, lockFileName),
		os.O_RDONLY,
		lockFilePerm)

	if os.IsNotExist(err) {
		return nil, ErrNoLock
	}

	if err != nil {
		return nil, err
	}

	var out LockFile

	if err := yaml.NewDecoder(file).Decode(&out); err != nil {
		return nil, err
	}

	if out.Modules == nil {
		out.Modules = make(map[string]*ModuleVersion)
	}

	return &out, file.Close()
}

// AddModule adds a new module entry to the lock file in the provided working directory.
func AddModule(l *LockFile, wd, name string) error {
	l.Modules[name] = &ModuleVersion{}

	return UpdateLockFile(wd, l)
}
