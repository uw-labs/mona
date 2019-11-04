package deps

import (
	"bytes"
	"io"
	"os/exec"
	"sort"
	"strings"

	"github.com/apex/log"
	"github.com/pkg/errors"
)

type AppDeps struct {
	Internal []string
	External []string
}

func GetAppDeps(mod Module, appPath string) (deps AppDeps, err error) {
	log.Debugf("Getting dependencies for app at %s", appPath)
	cmd := exec.Command("go", "list", "-f", `'{{ join .Deps "\n" }}'`, appPath)
	buf := &bytes.Buffer{}
	errBuf := &bytes.Buffer{}

	cmd.Stdout = buf
	cmd.Stderr = errBuf
	if err := cmd.Run(); err != nil {
		return deps, errors.Wrap(err, errBuf.String())
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
	sort.Sort(sort.StringSlice(deps.Internal))
	sort.Sort(sort.StringSlice(deps.External))

	return deps, nil
}
