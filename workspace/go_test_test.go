package main

import (
	"os/exec"
	"testing"

	. "github.com/smartystreets/assertions"
)

func TestNoGoFiles(t *testing.T) {
	output, err := GoTest("github.com/mdwhatcott/scantest/workspace/no_go_files")

	const expected = "can't load package: package github.com/mdwhatcott/scantest/workspace/no_go_files: no buildable Go source files "

	if ok, result := So(output, ShouldContainSubstring, expected); !ok {
		t.Error("\n" + result)
	}

	if ok, result := So(err, ShouldNotBeNil); !ok {
		t.Fatal(result)
	}

	if ok, result := So(err.Error(), ShouldEqual, "exit status 1"); !ok {
		t.Error("\n" + result)
	}
}

func TestNoTestFiles(t *testing.T) {
	output, err := GoTest("github.com/mdwhatcott/scantest/workspace/no_test_files")

	const expected = "?   	github.com/mdwhatcott/scantest/workspace/no_test_files	[no test files]"

	if ok, result := So(output, ShouldContainSubstring, expected); !ok {
		t.Error("\n" + result)
	}

	if ok, result := So(err, ShouldBeNil); !ok {
		t.Fatal(result)
	}
}

func TestNoTestFunctions(t *testing.T) {
	output, err := GoTest("github.com/mdwhatcott/scantest/workspace/no_test_functions")

	const expected = "ok  	github.com/mdwhatcott/scantest/workspace/no_test_functions	"

	if ok, result := So(output, ShouldContainSubstring, expected); !ok {
		t.Error(result)
	}

	if ok, result := So(err, ShouldBeNil); !ok {
		t.Fatal(result)
	}
}

func GoTest(folder string) (string, error) {
	command := exec.Command("go", "test", folder)
	output, err := command.CombinedOutput()
	return string(output), err
}
