// scantestgen continually scans the current directory for changes to
// _test.go files and runs go generate on the containing package.
// This is useful when working on packages tested using
// github.com/smartystreets/gunit.
package main

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
	"sync"
)

func main() {
	for {
		scanForChanges(root)
		goGenerate()
		reset()
	}
}
func scanForChanges(root string) {
	if err := filepath.Walk(root, scan); err != nil {
		log.Fatal(err)
	}

}

func scan(path string, info os.FileInfo, err error) error {
	if isGoTestFile(info) {
		fresh[filepath.Dir(path)] += checksum(info)
	}
	return nil
}

func isGoTestFile(file os.FileInfo) bool {
	return !file.IsDir() && strings.HasSuffix(file.Name(), "_test.go")
}

func checksum(file os.FileInfo) int64 {
	return file.Size() + file.ModTime().Unix()
}

func goGenerate() {
	for directory, state := range fresh {
		if stale[directory] != state {
			waiter.Add(1)
			go generate(directory)
		}
	}
}

func generate(path string) {
	pkg := packageName(path)
	log.Println("Running go generate for:", pkg)
	output, err := exec.Command("go", "generate", pkg).CombinedOutput()
	if err != nil {
		log.Printf("[ERROR] go generate %s:\n%s\n", pkg, string(output))
	}
	waiter.Done()
}

func packageName(path string) string {
	return path[len(goPath)+len("/src/"):]
}

func reset() {
	stale = fresh
	fresh = make(map[string]int64)
	time.Sleep(interval)
	waiter.Wait()
}

/**************************************************************************/

var (
	root   = resolveWorkingDirectory()
	stale  = make(map[string]int64)
	fresh  = make(map[string]int64)
	waiter = new(sync.WaitGroup)
	goPath = os.Getenv("GOPATH")
)

func resolveWorkingDirectory() string {
	root, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	return root
}

const interval = time.Millisecond * 250
