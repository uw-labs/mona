// Package environment contains types responsible for executing commands
// in different environments.
package environment

import "context"

type (
	// The Environment interface defines methods to be implemented for every
	// kind of environment we can execute commands in
	Environment interface {
		Execute(ctx context.Context, command string) error
	}
)
