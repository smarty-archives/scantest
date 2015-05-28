# scantest

A simple, repsonsive (a word which here means 'snappy') test runner. Like GoConvey, but smarter about what package to run, and with a much simpler interface, and at a fraction of the LOC.

## Features

- Runs `go test` for all packages under the current working directory.
- Scans for changes to .go files under the current directory.
- Runs tests for packages with changed .go files
- Runs tests for packages that depend on the modified package, if the change was not just in a _test.go file.
- Provides conventional output for console-based use or JSON for use with the command at github.com/smartystreets/scantest/scantest-web.
- Provides colorful output according to exit status of tests in both console and web mode (green=passed, red=failed).

### Installation

```
go get github.com/joewalnes/websocketd
go get github.com/smartystreets/scantest/...
```

### Console Runner Execution

```
cd my-project
scantest
```

Results of your tests will display in the terminal until you enter `<ctrl>+c`.

### Web Runner Execution

```
cd my-project
scantest-web
```

Then open your web browser to [`http://localhost:8888`](http://localhost:8888) to see your tests run. Save a change to a .go file somewhere under the current directory and see the tests for that package and any packages that depend on the modified package execute. Kill `scantest-web` by hitting `<ctrl>+c`.

## Custom Go Test Arguments

By default, `go test` is invoked with the `-v` flag and that's it. But sometimes you need control over how `go test` is invoked. Just plunk down a `.gotestargs` file in a package and when `go test` is run, the arguments listed there (on the first line) will be fed to `go test`.

Example contents:

```
-v -cover -parallel=5 -run=TestAllTheThings
```
