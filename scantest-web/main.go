package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

func main() {
	program := NewProgram()
	program.Prepare()
	program.Run()
}

type Program struct {
	command *exec.Cmd
	stdout  io.ReadCloser
}

func NewProgram() *Program {
	_, file, _, _ := runtime.Caller(0)
	static := filepath.Join(filepath.Dir(file), "../client")
	return &Program{
		command: exec.Command(
			"websocketd", "-port=8888", "-passenv=GOPATH,PATH", "--staticdir="+static, "scantest"),
	}
}

func (self *Program) Prepare() {
	self.command.Env = os.Environ()
	self.ConnectToStdout()
}
func (self *Program) Run() {
	self.command.Start()
	go self.PrintForevor()
	self.command.Wait()
}

func (self *Program) ConnectToStdout() {
	var err error
	self.stdout, err = self.command.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
}

func (self *Program) PrintForevor() {
	scanner := bufio.NewScanner(self.stdout)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
}
