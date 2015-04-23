
# Installation
```
go get github.com/joewalnes/websocketd
go get github.com/smartystreets/scantest
```

# Execution (Console-only)

```
cd my-project
scantest
```

# Execution (Web UI)

(This is a bit messy for now...)

```
websocketd -port=8080 -passenv=PATH,GOPATH --staticdir=$GOPATH/src/github.com/smartystreets/scantest/client scantest
```

Then open your web browser to [`http://localhost:8080`](http://localhost:8080) to see your tests run. Save a change to a .go file somewhere under the current directory and see the tests for that package and any packages that depend on the modified package execute.
