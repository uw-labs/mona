package git

import (
	"bytes"
	"os/exec"
	"strings"

	"github.com/apex/log"
	"github.com/pkg/errors"
)

// Branch returns the current git branch.
func Branch() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	log.Debug(cmd.String())

	buf := &bytes.Buffer{}
	errBuf := &bytes.Buffer{}

	cmd.Stdout = buf
	cmd.Stderr = errBuf
	if err := cmd.Run(); err != nil {
		return "", errors.Wrap(err, errBuf.String())
	}
	return strings.TrimSpace(buf.String()), nil
}
