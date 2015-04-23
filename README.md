
# Installation and Execution (Console Runner only)

```
go get github.com/smartystreets/scantest
cd my-project
scantest
```

Results of your tests will display in the terminal until you enter `<ctrl>+c`.

# Installation and Execution (Web Runner and/or Console Runner)

```
go get github.com/joewalnes/websocketd
go get github.com/smartystreets/scantest
go get github.com/smartystreets/scantest-web
cd my-project
scantest-web
```

Then open your web browser to [`http://localhost:8888`](http://localhost:8888) to see your tests run. Save a change to a .go file somewhere under the current directory and see the tests for that package and any packages that depend on the modified package execute. Kill `scantest-web` by hitting `<ctrl>+c`.
