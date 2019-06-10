package output_test

import (
	"bytes"
	"testing"

	"github.com/davidsbond/mona/internal/output"
	"github.com/stretchr/testify/assert"
)

func TestWritef(t *testing.T) {
	t.Parallel()

	tt := []struct {
		Name     string
		Format   string
		Args     []interface{}
		Expected string
	}{
		{
			Name:     "It should write a line with formatted arguments",
			Format:   "My %s is %s",
			Args:     []interface{}{"string", "formatted"},
			Expected: "My string is formatted\n",
		},
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			buff := bytes.NewBuffer([]byte{})

			if err := output.Writef(buff, tc.Format, tc.Args...); err != nil {
				assert.Fail(t, err.Error())
				return
			}

			assert.Equal(t, tc.Expected, buff.String())
		})
	}
}
