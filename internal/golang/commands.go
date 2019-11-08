package golang

import (
	"io"
	"os/exec"
	"strings"

	"github.com/uw-labs/mona/internal/executil"
)

func FindAllCommands(modName string) ([]string, error) {
	cmd := exec.Command("go", "list", "./...")
	buf, err := executil.RunCommand(cmd)
	if err != nil {
		return nil, err
	}

	allCMDs := make([]string, 0)
	for {
		line, err := buf.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		line = strings.TrimSpace(line)
		parts := strings.Split(line, "/")
		if len(parts) > 1 && parts[len(parts)-2] == "cmd" {
			line := strings.TrimPrefix(line, modName)
			allCMDs = append(allCMDs, strings.TrimPrefix(line, "/"))
		}
	}
	return allCMDs, nil
}
