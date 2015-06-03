package main

import (
	"flag"
	"strings"

	"github.com/smartystreets/scantest/scantest/contract"
	"github.com/smartystreets/scantest/scantest/shell"
)

func parseConfiguration() Config {
	root := flag.String("root", ".", "Specifies the root directory to scan for tests.")
	ignoredFolders := flag.String("ignore-dirs", ".git,.hg", "Comma-delimited list of folder names whose contents should be ignored.")
	flag.Parse()

	config := Config{
		RootFolder:     *root,
		IgnoredFolders: strings.Split(*ignoredFolders, ","),
	}
	return config
}

type Config struct {
	RootFolder     string
	IgnoredFolders []string
}

func buildHandlers(config Config) []contract.Handler {
	return []contract.Handler{
		shell.NewFileSystemScanner(config.RootFolder, config.IgnoredFolders),
	}
}

/*
   file system scanner
   checksummer
   import modified package
   collect upstream packages
   package sorter
   go generator
   gunit directive validator
   test args reader
   test runner
   test result interpreter
   test failure parser
   result printer
*/
