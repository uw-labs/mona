package output

import (
	"bytes"
	"fmt"
	"io"
)

// Writef formats a string with the given argument and writes it to the io.Writer
// implementation, adding a newline at the end.
func Writef(out io.Writer, format string, args ...interface{}) error {
	line := fmt.Sprintf(format, args...)
	buff := bytes.NewBufferString(line)
	buff.WriteRune('\n')

	_, err := out.Write(buff.Bytes())
	return err
}
