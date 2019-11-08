package golang

import (
	"bufio"
	"io"
	"os"
	"strings"
)

type Module struct {
	Name string
	// Requires is a map of dependency to its version
	Requires map[string]string
}

func ParseModuleFile(moduleFile string) (mod Module, err error) {
	f, err := os.Open(moduleFile)
	if err != nil {
		return mod, err
	}
	defer f.Close()

	return ParseModule(f)
}

func ParseModule(reader io.Reader) (mod Module, err error) {
	r := bufio.NewReader(reader)

	// Read module name
	line, err := r.ReadString('\n')
	if err != nil {
		return mod, err
	}
	mod.Name = strings.TrimPrefix(strings.TrimSpace(line), "module ")
	mod.Requires = make(map[string]string)

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

		if readingDeps {
			line = strings.TrimSpace(line)
			line = strings.TrimSuffix(line, "// indirect")
			parts := strings.Split(line, " ")
			mod.Requires[parts[0]] = parts[1]
			// TODO: need to deal with replace and edit,
			//  might be better to use some go tooling
		}
	}
}
