package config

import (
	"fmt"
)

type CommandConfig struct {
	Env   map[string]string `yaml:"env"`
	Flags []string          `yaml:"flags"`
}

// EnvToList converts the Env map to list of entries in form key=value.
func (cc CommandConfig) EnvToList() []string {
	result := make([]string, 0, len(cc.Env))
	for k, v := range cc.Env {
		result = append(result, fmt.Sprintf("%s=%s", k, v))
	}
	return result
}
