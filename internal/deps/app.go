package deps

import (
	"bytes"
	"io"
	"os/exec"
	"strings"
)

type Deps struct {
	Internal []string
	External []string
}

func GetForApp(mod Module, appPath string) (deps Deps, err error) {
	cmd := exec.Command("go", "list", "-f", `'{{ join .Deps "\n" }}'`, appPath)
	buf := &bytes.Buffer{}

	cmd.Stdout = buf
	if err := cmd.Run(); err != nil {
		return deps, err
	}

	allDeps := make([]string, 0)
	for {
		line, err := buf.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return deps, err
		}
		allDeps = append(allDeps, strings.TrimSpace(line))
	}
	seen := make(map[string]bool)
	for _, dep := range allDeps {
		if strings.HasPrefix(dep, mod.Name) {
			deps.Internal = append(deps.Internal, strings.TrimPrefix(dep, mod.Name+"/"))
			continue
		}
		for m, v := range mod.Deps {
			if strings.HasPrefix(dep, m) && !seen[m] {
				deps.External = append(deps.External, m+"@"+v)
				seen[m] = true
				break
			}
		}
	}
	return deps, nil
}
