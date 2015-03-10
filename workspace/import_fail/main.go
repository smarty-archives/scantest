package main

import (
	"fmt"
	"go/build"
)

func main() {
	p, err := build.ImportDir("/Users/mike/src/github.com/mdwhatcott/scantest/workspace/import_fail/fail", build.AllowBinary)
	if err != nil {
		fmt.Println("ERROR:", err)
	} else {
		fmt.Println(p.Name)
		fmt.Println(p.ImportPath)
		fmt.Println(p.IsCommand())
		fmt.Printf("%#v", p)
	}
}
