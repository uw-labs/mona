// Package config contains global configuration values used by the application.
package config

var (
	// Parallelism is used to determine the number of goroutines to use for performing
	// concurrent file walking. Increasing this value can improve scanning times in
	// large repositories. Is 10 by default.
	Parallelism = 10
)
