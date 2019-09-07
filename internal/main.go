package main

import (
	"fmt"

	"github.com/davidsbond/mona/internal/deps"
)

func main() {
	m, err := deps.ParseModule("../go.mod")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Module:", m.Name)
	for k, v := range m.Deps {
		fmt.Println("\t", k, "=>", v)
	}

	deps, err := deps.GetForApp(m, "../")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Internal:")
	for _, v := range deps.Internal {
		fmt.Println("\t", v)
	}
	fmt.Println("External:")
	for _, v := range deps.External {
		fmt.Println("\t", v)
	}
}
