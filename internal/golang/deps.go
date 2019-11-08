package golang

import (
	"io"
	"os/exec"
	"sort"
	"strings"

	"github.com/apex/log"

	"github.com/uw-labs/mona/internal/executil"
)

type Dependencies struct {
	Internal []string
	External []string
}

func GetDependencies(pkg string, mod Module) (deps Dependencies, err error) {
	log.Debugf("Getting dependencies for %s", pkg)
	cmd := exec.Command("go", "list", "-f", `'{{ join .Deps "\n" }}'`, "./"+pkg)
	buf, err := executil.RunCommand(cmd)
	if err != nil {
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
		for m, v := range mod.Requires {
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
