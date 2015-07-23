package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/smartystreets/scantest/go-shlex"
)

func main() {
	command := flag.String("command", deriveDefaultCommand(), "The command (with arguments) to run when a .go file is saved.")
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

// deriveDefaultCommand determines what the present as the default value of the
// command flag, in case the user does not provide one. It first looks for a
// Makefile in the current directory. If that doesn't exist it looks for the
// Makefile provided by this project, which serves as a working generic example
// that should fit a variety of use cases. If that doesn't exist for whatever
// reason (say, if the scantest binary wasn't built from source on the current
// machine) then it defaults to 'go test'.
func deriveDefaultCommand() string {
	var defaultCommand string

	if current, err := os.Getwd(); err == nil {
		if _, err := os.Stat(filepath.Join(current, "Makefile")); err == nil {
			defaultCommand = "make"
		}
	}

	if defaultCommand == "" {
		if _, file, _, ok := runtime.Caller(0); ok {
			backupMakefile := filepath.Join(filepath.Dir(file), "Makefile")
			if _, err := os.Stat(backupMakefile); err == nil {
				defaultCommand = backupMakefile
			}
		}
	}

	if defaultCommand == "" {
		defaultCommand = "go test"
	}

	return defaultCommand
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
	fmt.Fprintf(os.Stdout, clearScreen)
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
	fmt.Fprintln(os.Stdout, string(output))
	fmt.Fprintln(os.Stdout, strings.Repeat("-", len(message)))
	fmt.Fprintf(os.Stdout, resetColor)
}

func (this *Runner) run() (output []byte, success bool) {
	command := exec.Command(this.command[0])
	if len(this.command) > 1 {
		command.Args = this.command[1:]
	}
	command.Dir = this.working

	until := make(chan bool, 1)
	now := time.Now()
	go spin(now, until)

	var err error
	output, err = command.CombinedOutput()
	until <- true
	fmt.Println(Round(time.Since(now), time.Millisecond))
	if err != nil {
		output = append(output, []byte(err.Error())...)
	}
	return output, command.ProcessState.Success()
}

func spin(now time.Time, finished chan bool) {
	time.Sleep(time.Millisecond * 500)
	for {
		select {
		case <-finished:
			return
		default:
			fmt.Println(Round(time.Since(now), time.Millisecond))
			time.Sleep(time.Millisecond * 500)
		}
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
	clearScreen = "\033[2J\033[H" // clear the screen and put the cursor at top-left
	greenColor  = "\033[32m"
	redColor    = "\033[31m"
	resetColor  = "\033[0m"
)
