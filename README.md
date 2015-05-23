# scantest

A simple, repsonsive (a word which here means 'snappy') test runner. Like GoConvey, but smarter about what package to run, and with a much simpler interface, and at a fraction of the LOC.

## Features

- Runs `go test` for all packages under the current working directory.
- Scans for changes to .go files under the current directory.
- Runs tests for packages with changed .go files
- Runs tests for packages that depend on the modified package, if the change was not just in a _test.go file.
- Provides conventional output for console-based use or JSON for use with the command at github.com/smartystreets/scantest/scantest-web.
- Provides colorful output according to exit status of tests in both console and web mode (green=passed, red=failed).

### Installation and Execution (Console Runner only)

```
go get github.com/smartystreets/scantest
cd my-project
scantest
```

Results of your tests will display in the terminal until you enter `<ctrl>+c`.

### Installation and Execution (Web Runner and/or Console Runner)

```
go get github.com/joewalnes/websocketd
go get github.com/smartystreets/scantest
go get github.com/smartystreets/scanatest/scantest-web
cd my-project
scantest-web
```

Then open your web browser to [`http://localhost:8888`](http://localhost:8888) to see your tests run. Save a change to a .go file somewhere under the current directory and see the tests for that package and any packages that depend on the modified package execute. Kill `scantest-web` by hitting `<ctrl>+c`.
