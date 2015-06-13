package core

import (
	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
	"github.com/smartystreets/scantest/scantest/contract"
)

type ChecksummerFixture struct {
	*gunit.Fixture
	checksummer *Checksummer
	context     *contract.Context
}

func (this *ChecksummerFixture) Setup() {
	this.checksummer = NewChecksummer()
	this.context = &contract.Context{}
}

func (this *ChecksummerFixture) Test_NoFiles_ContextComplete() {
	this.context.Files = nil
	this.checksummer.Handle(this.context)
	this.So(this.context.Error, should.Equal, contract.ContextComplete)
}

func (this *ChecksummerFixture) Test_NoModifiedFiles_ContextComplete() {
	for x := 0; x < 10; x++ {
		this.addStaleGoFile(42)
	}
	this.checksummer.Handle(this.context)
	this.So(this.context.Error, should.Equal, contract.ContextComplete)
}

func (this *ChecksummerFixture) Test_ModifiedGoFile_ContextContinues() {
	this.addGoFile(42, 43)
	this.checksummer.Handle(this.context)
	this.So(this.context.Error, should.BeNil)
}

func (this *ChecksummerFixture) Test_SecondPassNoFilesModified_ContextComplete() {
	this.addGoFile(42, 43)
	this.checksummer.Handle(this.context)
	this.checksummer.Handle(this.context)
	this.So(this.context.Error, should.Equal, contract.ContextComplete)
}

func (this *ChecksummerFixture) Test_GenereatedFilesShouldNotBeSummed() {
	this.addGeneratedFiles()
	this.checksummer.Handle(this.context)
	this.So(this.context.Error, should.Equal, contract.ContextComplete)
}

func (this *ChecksummerFixture) Test_ArgumentFilesShouldBeSummed() {
	this.addArgumentFile()
	this.checksummer.Handle(this.context)
	this.So(this.context.Error, should.BeNil)
}

func (this *ChecksummerFixture) addFile(file *contract.File) {
	this.context.Files = append(this.context.Files, file)
}
func (this *ChecksummerFixture) addStaleGoFile(size int64) {
	this.addFile(&contract.File{
		Size:     size,
		IsGoFile: true,
	})
}
func (this *ChecksummerFixture) addGoFile(size, modified int64) {
	this.addFile(&contract.File{
		IsGoFile:   true,
		IsModified: true,
		Size:       size,
		Modified:   modified,
	})
}
func (this *ChecksummerFixture) addGeneratedFiles() {
	this.addFile(&contract.File{
		Path:       "/hello/generated_by_gunit_test.go",
		IsGoFile:   true,
		IsModified: true,
		Size:       42,
		Modified:   45,
	})
}
func (this *ChecksummerFixture) addArgumentFile() {
	this.addFile(&contract.File{
		Path:       "/hello/.testargs",
		IsModified: true,
		Size:       42,
		Modified:   45,
	})
}
