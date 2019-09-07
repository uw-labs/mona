package deps

import (
	"bufio"
	"io"
	"os"
	"strings"
)

type Module struct {
	Name string
	// Deps is a map of dependency to its version
	Deps map[string]string
}

func ParseModule(moduleFile string) (mod Module, err error) {
	f, err := os.Open(moduleFile)
	if err != nil {
		return mod, err
	}
	defer f.Close()

	r := bufio.NewReader(f)

	// Read module name
	line, err := r.ReadString('\n')
	if err != nil {
		return mod, err
	}
	mod.Name = strings.TrimPrefix(strings.TrimSpace(line), "module ")
	mod.Deps = make(map[string]string)

	// Read dependencies
	readingDeps := false
	for {
		line, err = r.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				return mod, nil
			}
			return mod, err
		}
		if line == "require (\n" || line == ")\n" {
			// Star stop processing dependencies
			readingDeps = !readingDeps
			continue
		}

		if readingDeps && !strings.Contains(line, "// indirect") {
			// We ignore indirect dependencies
			line = strings.TrimSpace(line)
			parts := strings.Split(line, " ")
			mod.Deps[parts[0]] = parts[1]
			// TODO: need to deal with replace and edit,
			//  might be better to use some go tooling
		}
	}
}
