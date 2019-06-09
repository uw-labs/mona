package output_test

import (
	"bytes"
	"testing"

	"github.com/davidsbond/mona/internal/output"
	"github.com/stretchr/testify/assert"
)

func TestWriteList(t *testing.T) {
	t.Parallel()

	tt := []struct {
		Name     string
		Title    string
		Items    []string
		Expected string
	}{
		{
			Name:     "It should format a list of items",
			Title:    "A test list",
			Items:    []string{"1", "2", "3"},
			Expected: "A test list\n- 1\n- 2\n- 3\n\n",
		},
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			buff := bytes.NewBuffer([]byte{})

			if err := output.WriteList(buff, tc.Title, tc.Items); err != nil {
				assert.Fail(t, err.Error())
				return
			}

			assert.Equal(t, tc.Expected, buff.String())
		})
	}
}
