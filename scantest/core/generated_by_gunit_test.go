//////////////////////////////////////////////////////////////////////////////
// Generated Code ////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////

package core

import (
	"testing"

	"github.com/smartystreets/gunit"
)

///////////////////////////////////////////////////////////////////////////////

func TestChecksummerFixture(t *testing.T) {
	fixture := gunit.NewFixture(t)

	test0 := &ChecksummerFixture{Fixture: fixture}
	test0.RunTestCase__(test0.Test_NoFiles_ContextComplete, "No files context complete")

	test1 := &ChecksummerFixture{Fixture: fixture}
	test1.RunTestCase__(test1.Test_NoModifiedFiles_ContextComplete, "No modified files context complete")

	test2 := &ChecksummerFixture{Fixture: fixture}
	test2.RunTestCase__(test2.Test_ModifiedGoFile_ContextContinues, "Modified go file context continues")

	test3 := &ChecksummerFixture{Fixture: fixture}
	test3.RunTestCase__(test3.Test_SecondPassNoFilesModified_ContextComplete, "Second pass no files modified context complete")

	test4 := &ChecksummerFixture{Fixture: fixture}
	test4.RunTestCase__(test4.Test_GenereatedFilesShouldNotBeSummed, "Genereated files should not be summed")

	test5 := &ChecksummerFixture{Fixture: fixture}
	test5.RunTestCase__(test5.Test_ArgumentFilesShouldBeSummed, "Argument files should be summed")

	fixture.Finalize()
}

func (self *ChecksummerFixture) RunTestCase__(test func(), description string) {
	self.Describe(description)
	self.Setup()
	test()
}

///////////////////////////////////////////////////////////////////////////////

func init() {
	gunit.Validate("0b197993c82d1d6441c89b1202b95f3b")
}

///////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////// Generated Code //
///////////////////////////////////////////////////////////////////////////////
