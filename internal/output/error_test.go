package output_test

import (
	"bytes"
	"errors"
	"testing"

	"github.com/davidsbond/mona/internal/output"
	"github.com/stretchr/testify/assert"
)

func TestWriteError(t *testing.T) {
	tt := []struct {
		Name     string
		Error    error
		Expected string
	}{
		{
			Name:     "It should format error messages",
			Error:    errors.New("a test error"),
			Expected: "A test error\n",
		},
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			buff := bytes.NewBuffer([]byte{})

			if err := output.WriteError(buff, tc.Error); err != nil {
				assert.Fail(t, err.Error())
				return
			}

			assert.Equal(t, tc.Expected, buff.String())
		})
	}
}
