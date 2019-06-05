package output

import (
	"bytes"
	"io"
	"unicode"
)

// WriteError writes a given error to the provider io.Writer implementation and
// ensures the first letter is upper-case. This means in code we can have lower-case
// errors to follow the standard, but when presented to a user they're formatted.
func WriteError(out io.Writer, err error) error {
	str := err.Error()

	for i, v := range str {
		// Convert the first character to upper case
		str = string(unicode.ToUpper(v)) + str[i+1:]
		break
	}

	builder := bytes.NewBufferString(str)
	builder.WriteRune('\n')

	_, err = io.Copy(out, builder)
	return err
}
