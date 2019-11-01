package config

import (
	"fmt"
	"strings"
)

type BuildConfig struct {
	Env     map[string]string `yaml:"env"`
	Flags   []string          `yaml:"flags"`
	LDFlags struct {
		S bool              `yaml:"s"`
		X map[string]string `yaml:"X"`
	} `yaml:"ldflags"`
}

func (bc BuildConfig) AllFlags() []string {
	ldFlags := make([]string, 0)
	if bc.LDFlags.S {
		ldFlags = append(ldFlags, "-s")
	}
	for k, v := range bc.LDFlags.X {
		ldFlags = append(ldFlags, fmt.Sprintf(`-X "%s=%s"`, k, v))
	}

	return append(
		bc.Flags,
		fmt.Sprintf("-ldflags=%s", strings.Join(ldFlags, " ")),
	)
}

// EnvToList converts the Env map to list of entries in form key=value.
func (bc BuildConfig) EnvToList() []string {
	return envToList(bc.Env)
}

type CommandConfig struct {
	Exclude    []string          `yaml:"exclude"`
	ExcludeMap map[string]bool   `yaml:"-"`
	Env        map[string]string `yaml:"env"`
	Flags      []string          `yaml:"flags"`
}

// EnvToList converts the Env map to list of entries in form key=value.
func (cc CommandConfig) EnvToList() []string {
	return envToList(cc.Env)
}

func envToList(env map[string]string) []string {
	result := make([]string, 0, len(env))
	for k, v := range env {
		result = append(result, fmt.Sprintf("%s=%s", k, v))
	}
	return result
}
