package shell

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/smartystreets/scantest/scantest/contract"
)

type FileSystemScanner struct {
	root           string
	ignoredFolders []string

	context *contract.Context
}

func NewFileSystemScanner(root string, ignoredFolders []string) *FileSystemScanner {
	return &FileSystemScanner{root: root}
}

func (this *FileSystemScanner) Handle(context *contract.Context) {
	this.context = context
	this.context.Error = filepath.Walk(this.root, this.walk)
}

func (this *FileSystemScanner) walk(path string, info os.FileInfo, err error) error {
	if this.isIgnoredFolder(info) {
		return filepath.SkipDir
	}

	// TODO: this should be done by the checksummer.
	// if info.Name() == "generated_by_gunit_test.go" {
	// 	return nil
	// }

	this.context.Files = append(this.context.Files, &contract.File{
		Path:         path,
		ParentFolder: filepath.Dir(path), // does this get the parent of a dir?
		IsFolder:     info.IsDir(),
		Size:         info.Size(),
		Modified:     info.ModTime().Unix(),
		IsGoFile:     strings.HasSuffix(path, ".go"),
		IsGoTestFile: strings.HasSuffix(path, "_test.go"),
	})

	return nil
}

func (this *FileSystemScanner) isIgnoredFolder(info os.FileInfo) bool {
	if !info.IsDir() {
		return false
	}
	for _, ignoredFolder := range this.ignoredFolders {
		if info.Name() == ignoredFolder {
			return true
		}
	}
	return false
}
