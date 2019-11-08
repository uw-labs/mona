package git

import (
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"

	"github.com/uw-labs/mona/internal/golang"
)

type GoDiff struct {
	Packages map[string]bool
	Modules  map[string]bool
}

func (diff GoDiff) Changed(pkg string, mod golang.Module) (bool, error) {
	if diff.Packages[pkg] {
		return true, nil
	}
	deps, err := golang.GetDependencies(pkg, mod)
	if err != nil {
		return false, err
	}
	for _, dep := range deps.Internal {
		if diff.Packages[dep] {
			return true, nil
		}
	}
	for _, mod := range deps.External {
		if diff.Modules[mod] {
			return true, nil
		}
	}
	return false, nil
}

func (diff GoDiff) PackagesList() []string {
	out := make([]string, 0, len(diff.Packages))
	for pkg := range diff.Packages {
		out = append(out, pkg)
	}
	return out
}

func GetGoDiff(goModAfter golang.Module, branch string) (GoDiff, error) {
	diff := GoDiff{
		Packages: make(map[string]bool),
		Modules:  make(map[string]bool),
	}
	cmd := exec.Command("git", "diff", "--name-only", branch)
	buf := &bytes.Buffer{}
	errBuf := &bytes.Buffer{}

	cmd.Stdout = buf
	cmd.Stderr = errBuf
	if err := cmd.Run(); err != nil {
		return diff, errors.Wrap(err, errBuf.String())
	}

	goModChanged := false
	for {
		line, err := buf.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return diff, err
		}
		line = strings.TrimSpace(line)

		if !goModChanged && line == "go.mod" {
			goModChanged = true
			continue
		}
		if strings.HasSuffix(line, ".go") {
			pkg, _ := filepath.Split(line)
			pkg = strings.TrimSuffix(pkg, "/")
			diff.Packages[pkg] = true
		}
	}
	if goModChanged {
		goModBefore, err := getGoModBefore(branch)
		if err != nil {
			return diff, err
		}
		// We can ignore any dependencies that were removed as they must have been
		// removed from a file that used them, so this package will be checked anyway.
		for mod, version := range goModAfter.Requires {
			oldVersion, ok := goModBefore.Requires[mod]
			if !ok {
				// This dependency wasn't used before so there must be
				// a changed file or files that imports so we can ignore it,
				// as those will be checked regardless.
				continue
			}
			if oldVersion != version {
				diff.Modules[mod] = true
			}
		}
	}

	return diff, nil
}

func getGoModBefore(branch string) (mod golang.Module, err error) {
	cmd := exec.Command("git", "show", fmt.Sprintf("%s:go.mod", branch))
	buf := &bytes.Buffer{}
	errBuf := &bytes.Buffer{}

	cmd.Stdout = buf
	cmd.Stderr = errBuf
	if err := cmd.Run(); err != nil {
		return mod, errors.Wrap(err, errBuf.String())
	}
	return golang.ParseModule(buf)
}
