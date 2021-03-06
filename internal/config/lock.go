package config

import (
	"os"
	"path/filepath"

	"github.com/hashicorp/go-multierror"
	"gopkg.in/yaml.v2"
)

const (
	lockFileName = "mona.lock"
	lockFilePerm = 0644
)

type (
	// The LockFile type represents the structure of a lock file, it stores the project name,
	// version and the last build hashes used for each app
	LockFile struct {
		Apps map[string]*AppVersion `yaml:"apps,omitempty"`
	}

	// The AppVersion type represents individual app information as stored
	// in the lock file.
	AppVersion struct {
		BuildHash string `yaml:"build"`
		TestHash  string `yaml:"test"`
		LintHash  string `yaml:"lint"`
	}
)

// NewLockFile creates a new "mona.lock" file in the current working directory using the
// provided name.
func NewLockFile(dir string, name string) error {
	location := filepath.Join(dir, lockFileName)
	file, err := os.Create(location)

	if err != nil {
		return err
	}

	lock := LockFile{
		Apps: make(map[string]*AppVersion),
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
		os.O_CREATE|os.O_WRONLY,
		lockFilePerm)

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
		return &LockFile{
			Apps: make(map[string]*AppVersion),
		}, nil
	}

	if err != nil {
		return nil, err
	}

	var out LockFile

	if err := yaml.NewDecoder(file).Decode(&out); err != nil {
		return nil, err
	}

	if out.Apps == nil {
		out.Apps = make(map[string]*AppVersion)
	}

	return &out, file.Close()
}

// AddApp adds a new app entry to the lock file in the provided working directory.
func AddApp(l *LockFile, wd, name string) error {
	l.Apps[name] = &AppVersion{}

	return UpdateLockFile(wd, l)
}
