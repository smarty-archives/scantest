package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/smartystreets/scantest/go-shlex"
)

func main() {
	command := flag.String("command", "make", "The command (with arguments) to run when a .go file is saved.")
	flag.Parse()

	working, err := os.Getwd()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	args, err := shlex.Split(*command)
	if err != nil {
		log.Fatal(err)
	}
	if len(args) < 1 {
		log.Fatal("Please provide something to run.")
	}

	scanner := &Scanner{working: working}
	runner := &Runner{working: working, command: args}
	for {
		if scanner.Scan() {
			runner.Run()
		}
	}
}

////////////////////////////////////////////////////////////////////////////

type Scanner struct {
	state   int64
	working string
}

func (this *Scanner) Scan() bool {
	time.Sleep(time.Millisecond * 250)
	newState := this.checksum()
	defer func() { this.state = newState }()
	return newState != this.state
}

func (this *Scanner) checksum() int64 {
	var sum int64 = 0
	err := filepath.Walk(this.working, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			sum++
		} else if info.Name() == "generated_by_gunit_test.go" {
			return nil
		} else if strings.HasSuffix(info.Name(), ".go") {
			sum += info.Size() + info.ModTime().Unix()
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	return sum
}

////////////////////////////////////////////////////////////////////////////

type Runner struct {
	command []string
	working string
}

func (this *Runner) Run() {
	message := fmt.Sprintln(" Executing:", strings.Join(this.command, " "))
	fmt.Fprintln(os.Stdout, "\n"+strings.Repeat("=", len(message)))
	fmt.Fprint(os.Stdout, message)
	fmt.Fprintln(os.Stdout, strings.Repeat("=", len(message)))
	output, problem := this.run()
	os.Stdout.Write(problem)
	fmt.Fprintln(os.Stdout)
	os.Stdout.Write(output)
}

func (this *Runner) run() (output []byte, errBytes []byte) {
	command := exec.Command(this.command[0])
	if len(this.command) > 1 {
		command.Args = this.command[1:]
	}
	command.Dir = this.working

	var err error
	output, err = command.CombinedOutput()
	if err != nil {
		errBytes = []byte(err.Error())
	}
	return output, errBytes
}

////////////////////////////////////////////////////////////////////////////
