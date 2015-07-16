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
	command  []string
	working  string
	finished bool
}

func (this *Runner) Run() {
	message := fmt.Sprintln(" Executing:", strings.Join(this.command, " "))
	fmt.Fprintln(os.Stdout, "\n"+strings.Repeat("=", len(message)))
	fmt.Fprint(os.Stdout, message)
	fmt.Fprintln(os.Stdout, strings.Repeat("=", len(message)))
	output, success := this.run()
	if success {
		fmt.Fprintf(os.Stdout, greenColor)
	} else {
		fmt.Fprintf(os.Stdout, redColor)
	}
	fmt.Fprintln(os.Stdout)
	os.Stdout.Write(output)
	fmt.Fprintln(os.Stdout, strings.Repeat("-", len(message)))
	fmt.Fprintf(os.Stdout, resetColor)
}

func (this *Runner) run() (output []byte, success bool) {
	command := exec.Command(this.command[0])
	if len(this.command) > 1 {
		command.Args = this.command[1:]
	}
	command.Dir = this.working

	this.finished = false
	go this.spin()

	var err error
	output, err = command.CombinedOutput()
	this.finished = true
	if err != nil {
		output = append(output, []byte(err.Error())...)
	}
	return output, command.ProcessState.Success()
}

func (this *Runner) spin() {
	now := time.Now()
	time.Sleep(time.Millisecond * 500)
	for !this.finished {
		fmt.Println(Round(time.Since(now), time.Millisecond))
		time.Sleep(time.Millisecond * 500)
	}
}

// GoLang-Nuts thread:
//     https://groups.google.com/d/msg/golang-nuts/OWHmTBu16nA/RQb4TvXUg1EJ
// Wise, a word which here means unhelpful, guidance from Commander Pike:
//     https://groups.google.com/d/msg/golang-nuts/OWHmTBu16nA/zoGNwDVKIqAJ
// Answer satisfying the original asker:
//     https://groups.google.com/d/msg/golang-nuts/OWHmTBu16nA/wnrz0tNXzngJ
// Answer implementation on the Go Playground:
//     http://play.golang.org/p/QHocTHl8iR
func Round(d, r time.Duration) time.Duration {
	if r <= 0 {
		return d
	}
	neg := d < 0
	if neg {
		d = -d
	}
	if m := d % r; m+m < r {
		d = d - m
	} else {
		d = d + r - m
	}
	if neg {
		return -d
	}
	return d
}

////////////////////////////////////////////////////////////////////////////

var (
	greenColor = "\033[32m"
	redColor   = "\033[31m"
	resetColor = "\033[0m"
)
