package main

import (
	"fmt"
	"go/build"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var (
	cascade       map[string][]string = make(map[string][]string)
	packageOfFile map[string]string   = make(map[string]string)
)

func main() {
	root := "/Users/mike/src/github.com/smartystreets/goconvey"

	// rootPackage, err := build.ImportDir(root, build.AllowBinary)
	// if err != nil {
	// 	log.Fatal("Couldn't import root package:", err)
	// }

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			return nil
		}

		_package, err := build.ImportDir(path, build.AllowBinary)

		if err != nil {
			return nil
		}

		if strings.Contains(path, root) {
			for _, _import := range append(_package.Imports, _package.TestImports...) {
				imported, err := build.Default.Import(_import, "", build.AllowBinary)
				if err != nil || imported.Goroot {
					continue
				}
				found := false
				for _, already := range cascade[_import] {
					if already == _package.ImportPath {
						found = true
					}
				}
				if !found {
					cascade[_import] = append(cascade[_import], _package.ImportPath)
				}
			}
		}
		return nil
	})

	if err != nil {
		log.Fatal("Problem gathering packages:", err)
	}

	fmt.Println()
	for up, downs := range cascade {
		fmt.Println(up)
		for _, down := range downs {
			fmt.Println("  " + down)
		}
		fmt.Println()
	}

	fmt.Println("\n--------\n")

	// for _, i := range downstream([]string{}, "github.com/smartystreets/goconvey/convey") {
	// 	fmt.Println(i)
	// }
}

func downstream(sofar []string, top string) []string {
	downs := cascade[top]
	if len(downs) == 0 {
		return sofar
	}

	sofar = append(sofar, downs...)

	for _, down := range downs {
		sofar = append(sofar, downstream(sofar, down)...)
	}
	return sofar
}
