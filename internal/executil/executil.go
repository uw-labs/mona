package executil

import (
	"bytes"
	"os/exec"

	"github.com/pkg/errors"
)

func RunCommand(cmd *exec.Cmd) (*bytes.Buffer, error) {
	errBuf := &bytes.Buffer{}
	resultBuf := &bytes.Buffer{}

	cmd.Stdout = resultBuf
	cmd.Stderr = errBuf
	if err := cmd.Run(); err != nil {
		return nil, errors.Wrap(err, errBuf.String())
	}
	return resultBuf, nil
}
