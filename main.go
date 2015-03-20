// go build && websocketd -port=8080 -passenv=PATH,GOPATH --staticdir=client ./scantest
// go install github.com/mdwhatcott/scantest && websocketd -port=8080 -passenv=PATH,GOPATH --staticdir=$GOPATH/src/github.com/mdwhatcott/scantest/client scantest
package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"go/build"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"syscall"
	"time"
)

//////////////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////////////

func main() {
	var pretty bool
	flag.BoolVar(&pretty, "pretty", false, "Set to true if you want pretty, multi-line output, or false if you want JSON (like for a browser).")
	flag.Parse()

	workingDirectory, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var (
		inputCommands = make(chan struct{})
		scannedFiles  = make(chan chan *File)
		checkedFiles  = make(chan chan *File)
		packages      = make(chan chan *Package)
		executions    = make(chan map[string]bool)
		results       = make(chan []Result)

		scanner = &FileSystemScanner{
			root: workingDirectory,
			out:  scannedFiles,
		}

		checksummer = &Checksummer{
			commands: inputCommands,

			in:  scannedFiles,
			out: checkedFiles,
		}

		packager = &Packager{
			in:  checkedFiles,
			out: packages,
		}

		selector = &PackageSelector{
			in:  packages,
			out: executions,
		}

		runner = &Runner{
			in:  executions,
			out: results,
		}

		printer = &Printer{
			in:     results,
			pretty: pretty,
		}
	)

	go scanner.ScanForever()
	go checksummer.RespondForevor()
	go checksummer.ListenForever()
	go packager.ListenForever()
	go selector.ListenForever()
	go runner.ListenForever()
	go printer.ListenForever()
	receiveInput(inputCommands)
}

//////////////////////////////////////////////////////////////////////////////////////

func receiveInput(signal chan struct{}) {
	for {
		a := []byte{0}
		os.Stdin.Read(a)
		if a[0] == 10 { // Enter key
			signal <- struct{}{}
		}
	}
}

//////////////////////////////////////////////////////////////////////////////////////

type File struct {
	Path         string
	ParentFolder string
	Size         int64
	Modified     int64
	IsFolder     bool
	IsGoFile     bool
	IsGoTestFile bool
	IsModified   bool
}

//////////////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////////////

type FileSystemScanner struct {
	root string
	out  chan chan *File
}

func (self *FileSystemScanner) ScanForever() {
	for {
		batch := make(chan *File)
		self.out <- batch

		filepath.Walk(self.root, func(path string, info os.FileInfo, err error) error { // TODO: handle err of filepath.Walk?
			if info.IsDir() && (info.Name() == ".git" || info.Name() == ".hg" /* etc... */) {
				return filepath.SkipDir
			}

			batch <- &File{
				Path:         path,
				ParentFolder: filepath.Dir(path), // does this get the parent of a dir?
				IsFolder:     info.IsDir(),
				Size:         info.Size(),
				Modified:     info.ModTime().Unix(),
				IsGoFile:     strings.HasSuffix(path, ".go"),
				IsGoTestFile: strings.HasSuffix(path, "_test.go"),
			}

			return nil
		})
		close(batch)
		time.Sleep(time.Millisecond * 250)
	}
}

//////////////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////////////

type Checksummer struct {
	commands chan struct{}
	reset    bool

	in  chan chan *File
	out chan chan *File

	state   int64
	goFiles map[string]int64
}

func (self *Checksummer) RespondForevor() {
	for {
		<-self.commands
		self.reset = true
		time.Sleep(time.Millisecond)
	}
}

func (self *Checksummer) ListenForever() {
	self.goFiles = map[string]int64{}

	for {
		state := int64(0)
		incoming := <-self.in
		outgoing := []*File{}
		goFiles := map[string]int64{}

		for file := range incoming {
			if !file.IsFolder && file.IsGoFile {
				fileChecksum := file.Size + file.Modified
				state += fileChecksum
				if checksum, found := self.goFiles[file.Path]; !found || checksum != fileChecksum {
					file.IsModified = true
				} else if self.reset { // the user has requested a re-run of all packages, so fake a modification.
					file.IsModified = true
				}
				goFiles[file.Path] = fileChecksum
				outgoing = append(outgoing, file)
			}
		}
		self.goFiles = goFiles

		if state != self.state || self.reset {
			fmt.Println("Running tests...")
			self.state = state
			out := make(chan *File)
			self.out <- out
			for _, file := range outgoing {
				out <- file
			}
			close(out)

			if self.reset {
				self.reset = false
			}
		}
	}
}

//////////////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////////////

type Package struct {
	Info           *build.Package
	IsModifiedTest bool
	IsModifiedCode bool
	// arguments string
}

//////////////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////////////

type Packager struct {
	in  chan chan *File
	out chan chan *Package
}

func (self *Packager) ListenForever() {
	for {
		incoming := <-self.in
		packages := map[string]*Package{} // key: Folder path

		for file := range incoming {
			pkg, found := packages[file.ParentFolder]
			if !found {
				pkg = &Package{}
				var err error
				pkg.Info, err = build.ImportDir(file.ParentFolder, build.AllowBinary)
				if err != nil {
					continue
				}
				packages[file.ParentFolder] = pkg
			}
			if file.IsModified && file.IsGoTestFile {
				pkg.IsModifiedTest = true
			} else if file.IsModified && !file.IsGoTestFile && file.IsGoFile {
				pkg.IsModifiedCode = true
			}
		}

		outgoing := make(chan *Package)
		self.out <- outgoing
		for _, pkg := range packages {
			outgoing <- pkg
		}
		close(outgoing)
	}
}

//////////////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////////////

type Execution struct {
	PackageName string
	// ParsedArguments []string
}

//////////////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////////////

type PackageSelector struct {
	in  chan chan *Package
	out chan map[string]bool
}

func (self *PackageSelector) ListenForever() {
	for {
		incoming := <-self.in
		executions := map[string]bool{}
		cascade := map[string][]string{}
		all := []*Package{}

		for pkg := range incoming {
			all = append(all, pkg)

			for _, _import := range append(pkg.Info.Imports, pkg.Info.TestImports...) {
				imported, err := build.Default.Import(_import, "", build.AllowBinary)
				if err != nil || imported.Goroot {
					continue
				}
				found := false
				for _, already := range cascade[_import] {
					if already == pkg.Info.ImportPath {
						found = true
					}
				}
				if !found {
					cascade[_import] = append(cascade[_import], pkg.Info.ImportPath)
				}
			}

			for _, pkg := range all {
				if pkg.IsModifiedCode || pkg.IsModifiedTest {
					executions[pkg.Info.ImportPath] = true
					if pkg.IsModifiedCode {
						for _, upstream := range cascade[pkg.Info.ImportPath] {
							executions[upstream] = true
						}
					}
				}
			}
		}

		self.out <- executions
	}
}

//////////////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////////////

type Result struct {
	PackageName string
	Status      PackageStatus
	Output      string
	Failures    []string
}

type PackageStatus int

const (
	CompileFailed PackageStatus = iota
	TestsFailed
	TestsPassed
)

//////////////////////////////////////////////////////////////////////////////////////

// ResultSet implements sort.Interface for []Person based on the result status and package name.
type ResultSet []Result

func (self ResultSet) Len() int      { return len(self) }
func (self ResultSet) Swap(i, j int) { self[i], self[j] = self[j], self[i] }
func (self ResultSet) Less(i, j int) bool {
	if self[i].Status == self[j].Status {
		return self[i].PackageName[0] < self[j].PackageName[0]
	}
	return self[i].Status < self[j].Status
}

//////////////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////////////

type Runner struct {
	in  chan map[string]bool
	out chan []Result
}

func (self *Runner) ListenForever() {
	for {
		results := []Result{}
		for packageName, _ := range <-self.in {
			command := exec.Command("go", "test", "-v", "-short", packageName) // TODO: profiles
			output, err := command.CombinedOutput()
			result := Result{
				PackageName: packageName,
				Output:      string(output),
			}

			// http://stackoverflow.com/questions/10385551/get-exit-code-go
			if err == nil { // if exit code is 0: the tests executed and passed.
				result.Status = TestsPassed
			} else if exit, ok := err.(*exec.ExitError); ok {
				if status, ok := exit.Sys().(syscall.WaitStatus); ok {

					if status.ExitStatus() == 1 { // if exit code is 1: we tests failed or panicked.
						result.Status = TestsFailed
						result.Failures = parseFailures(result)
					} else if status.ExitStatus() > 1 { // if exit code is > 1: we failed to build and tests were not run.
						result.Status = CompileFailed
					}
				}
			}

			results = append(results, result)
		}
		self.out <- results
	}
}

func parseFailures(result Result) []string {
	failures := []string{}
	if result.Status != TestsFailed {
		return failures
	}
	buffer := new(bytes.Buffer)
	reader := strings.NewReader(result.Output)
	scanner := bufio.NewScanner(reader)
	var passed bool

	for scanner.Scan() {
		line := scanner.Text() + "\n"
		if strings.HasPrefix(line, "=== RUN Test") {
			if buffer.Len() > 0 && passed == false {
				failures = append(failures, buffer.String())
			}
			buffer = bytes.NewBufferString("")
			buffer.WriteString(line)
		} else if strings.HasPrefix(line, "--- PASS: Test") {
			passed = true
		} else if strings.HasPrefix(line, "--- FAIL: Test") {
			buffer.WriteString(line)
			passed = false
		} else {
			buffer.WriteString(line)
		}
	}
	return failures
}

//////////////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////////////

type Printer struct {
	pretty bool
	in     chan []Result
}

func (self *Printer) ListenForever() {
	for resultSet := range self.in {
		sort.Sort(ResultSet(resultSet))
		if self.pretty {
			self.console(resultSet)
		} else {
			self.json(resultSet)
		}
	}
}

func (self *Printer) console(resultSet []Result) {
	const (
		red   = "\033[31m"
		green = "\033[32m"
		reset = "\033[0m"
	)

	failed := false

	for x := len(resultSet) - 1; x >= 0; x-- {
		result := resultSet[x]
		if result.Status < TestsPassed {
			failed = true
			fmt.Fprint(os.Stdout, red)
		}
		fmt.Println(result.PackageName)
		fmt.Println(result.Output)
		fmt.Println(reset)
		fmt.Println()
	}

	if failed {
		fmt.Fprint(os.Stdout, red)
	} else {
		fmt.Fprint(os.Stdout, green)
	}
	fmt.Println("-----------------------------------------------------")
	fmt.Println(reset)
}

type JSONResult struct {
	Packages []Result `json:"packages"`
}

func (self *Printer) json(resultSet []Result) {
	result := JSONResult{Packages: resultSet}
	raw, err := json.Marshal(result)
	if err != nil {
		fmt.Println(err)
		os.Exit(1) // TODO: maybe send a web socket message that indicates the UI of the crash...
	} else {
		fmt.Println(string(raw))
	}
}

//////////////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////////////
