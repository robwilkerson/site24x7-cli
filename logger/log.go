package logger

import (
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/pflag"
)

const (
	// SILENT should prevent any output from a command
	SILENT = -1
	WARN   = iota - 1
	INFO
	DEBUG
)

// SetVerbosity is a means of setting a persistent value to define output needs.
// Currently only one level of verbosity is supported.
func SetVerbosity(fs *pflag.FlagSet) int {
	quiet, _ := fs.GetBool("quiet")
	verbosity, _ := fs.GetCount("verbose")

	if quiet {
		verbosity = SILENT
	}

	if err := os.Setenv("VERBOSITY", fmt.Sprintf("%d", verbosity)); err != nil {
		panic(err.Error())
	}

	// fmt.Printf("--> Set verbosity = %d\n", verbosity)
	return verbosity
}

// GetVerbosity returns the command's verbosity setting
func GetVerbosity() int {
	if v, err := strconv.Atoi(os.Getenv("VERBOSITY")); err == nil {
		return v
	}

	// Default to no informational output
	return WARN
}

// Out writes a message to stdout as long as quiet mode is off; useful for
// command output
func Out(msg string) {
	if v := GetVerbosity(); v != SILENT {
		fmt.Fprintln(os.Stdout, msg)
	}
}

// Warn writes a key message to stderr so as not to interfere with handling json
// on stdout
func Warn(msg string) {
	if v := GetVerbosity(); v >= WARN {
		fmt.Fprintln(os.Stderr, "[ WARN ] ", msg)
	}
}

// Info writes an extended message to stderr so as not to interfere with handling
// json on stdout
func Info(msg string) {
	if v := GetVerbosity(); v >= INFO {
		fmt.Fprintln(os.Stderr, "[ INFO ] ", msg)
	}
}

// Debug writes a verbose message to stderr so as not to interfere with handling
// json on stdout
func Debug(msg string) {
	if v := GetVerbosity(); v >= DEBUG {
		fmt.Fprintln(os.Stderr, "[ DEBUG ]", msg)
	}
}
