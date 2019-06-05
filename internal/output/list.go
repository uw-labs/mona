package output

import (
	"bytes"
	"io"
)

// WriteList generates a text based list of the given items underneath the given
// title and writes it to the provided writer implementation.
func WriteList(out io.Writer, title string, items []string) error {
	if len(items) == 0 {
		return nil
	}

	builder := bytes.NewBuffer([]byte{})
	builder.WriteString(title)
	builder.WriteRune('\n')

	for _, item := range items {
		builder.WriteString("- ")
		builder.WriteString(item)
		builder.WriteRune('\n')
		builder.WriteRune('\n')
	}

	_, err := io.Copy(out, builder)
	return err
}
