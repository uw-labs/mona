package files

import (
	"errors"
	"fmt"
	"os"
	"strings"

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
		Name    string   `yaml:"name"`
		Version string   `yaml:"version"`
		Modules []string `yaml:"modules,omitempty"`
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
	info := fmt.Sprintf("%s %s %s", name, location, hash)
	l.Modules = append(l.Modules, info)

	file, err := os.Create(lockFileName)

	if err != nil {
		return err
	}

	defer file.Close()
	return yaml.NewEncoder(file).Encode(l)
}

// ParseLockLine reads a module line from the "mona.lock" file and splits it into
// its name, location and hash.
func ParseLockLine(line string) (name, location, hash string) {
	parts := strings.Split(line, " ")

	name = parts[0]
	location = parts[1]
	hash = parts[2]

	return
}

// CreateLockLine formats a new module line for the provided name,
// location and hash.
func CreateLockLine(name, location, hash string) string {
	return fmt.Sprintf("%s %s %s", name, location, hash)
}
