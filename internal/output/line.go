package output

import (
	"bytes"
	"fmt"
	"io"
)

func Writef(out io.Writer, format string, args ...interface{}) error {
	line := fmt.Sprintf(format, args...)
	buff := bytes.NewBufferString(line)
	buff.WriteRune('\n')

	_, err := out.Write(buff.Bytes())
	return err
}
