package contract

import "errors"

type Context struct {
	Files         []*File
	Checksum      int64
	ModifiedFiles []File
	Packages      []Package

	Error error
}

var ContextComplete = errors.New("Context complete")

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

type Package struct {
	ImportPath    string
	Imports       []string
	TestImports   []string
	GenerateError error
	TestArguments []string
	TestOutput    string
	TestError     error
	TestExitCode  int
}

// PackageList implements sort.Interface for []Package based on the result status and package name.
type PackageList []Package

func (self PackageList) Len() int      { return len(self) }
func (self PackageList) Swap(i, j int) { self[i], self[j] = self[j], self[i] }
func (self PackageList) Less(i, j int) bool {
	// if self[i].Status == self[j].Status {
	return self[i].ImportPath[0] < self[j].ImportPath[0]
	// }
	// return self[i].Status < self[j].Status
}
