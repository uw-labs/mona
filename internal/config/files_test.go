package config_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func deleteLockFile(t *testing.T) {
	if err := os.Remove("mona.lock"); err != nil {
		assert.Fail(t, err.Error())
	}
}

func deleteProjectFile(t *testing.T) {
	if err := os.Remove("mona.yml"); err != nil {
		assert.Fail(t, err.Error())
	}
}
