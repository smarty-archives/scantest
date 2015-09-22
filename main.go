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

	"github.com/mdwhatcott/spin"
	"github.com/smartystreets/scantest/go-shlex"
)

func main() {
	command := flag.String("command", deriveDefaultCommand(), "The command (with arguments) to run when a .go file is saved.")
	flag.Parse()

	working, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
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

// deriveDefaultCommand determines what to present as the default value of the
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
	write(".")
	return newState != this.state
}

func (this *Scanner) checksum() int64 {
	var sum int64 = 0
	err := filepath.Walk(this.working, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			sum++
		} else if info.Name() == "generated_by_gunit_test.go" {
			return nil
		} else if strings.HasSuffix(info.Name(), ".go") || info.Name() == "Makefile" {
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

	// write(clearScreen)
	writeln()
	write(strings.Repeat("=", len(message)))
	writeln()
	write(message)
	write(strings.Repeat("=", len(message)))
	writeln()
	output, success := this.run()
	if success {
		write(greenColor)
	} else {
		write(redColor)
	}
	write(string(output))
	writeln()
	write(strings.Repeat("-", len(message)))
	writeln()
	write(resetColor)
}

func writeln() {
	write("\n")
}
func write(a ...interface{}) {
	fmt.Fprint(os.Stdout, a...)
	os.Stdout.Sync()
}

func (this *Runner) run() (output []byte, success bool) {
	command := exec.Command(this.command[0])
	if len(this.command) > 1 {
		command.Args = append(command.Args, this.command[1:]...)
	}
	command.Dir = this.working

	now := time.Now()
	spinner := spin.New(spin.StyleLine, time.Millisecond*100)
	spinner.GoStart()

	var err error
	output, err = command.CombinedOutput()
	spinner.Stop()
	fmt.Println(Round(time.Since(now), time.Millisecond))
	if err != nil {
		output = append(output, []byte(err.Error())...)
	}
	return output, command.ProcessState.Success()
}

// GoLang-Nuts thread:
//     https://groups.google.com/d/msg/golang-nuts/OWHmTBu16nA/RQb4TvXUg1EJ
// Wise, a word which here means unhelpful, guidance from Commander Pike:
//     https://groups.google.com/d/msg/golang-nuts/OWHmTBu16nA/zoGNwDVKIqAJ
// Answer satisfying the original asker:
//     https://groups.google.com/d/msg/golang-nuts/OWHmTBu16nA/wnrz0tNXzngJ
// Answer implementation on the Go Playground:
//     http://play.golang.org/p/QHocTHl8iR
func Round(duration, precision time.Duration) time.Duration {
	if precision <= 0 {
		return duration
	}
	negative := duration < 0
	if negative {
		duration = -duration
	}
	if m := duration % precision; m+m < precision {
		duration = duration - m
	} else {
		duration = duration + precision - m
	}
	if negative {
		return -duration
	}
	return duration
}

////////////////////////////////////////////////////////////////////////////

var (
	clearScreen = "\033[2J\033[H" // clear the screen and put the cursor at top-left
	greenColor  = "\033[32m"
	redColor    = "\033[31m"
	resetColor  = "\033[0m"
)
