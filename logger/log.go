package logger

import (
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/pflag"
)

const (
	silent = -1
	warn   = iota - 1
	info
	debug
)

// SetLevel is a means of setting a persistent value to define output needs.
// Currently only one level of verbosity is supported.
func SetVerbosity(fs *pflag.FlagSet) int {
	quiet, _ := fs.GetBool("quiet")
	verbosity, _ := fs.GetCount("verbose")

	if quiet {
		verbosity = silent
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
	return warn
}

// Out writes a message to stdout as long as quiet mode is off; useful for
// command output
func Out(msg string) {
	if v := GetVerbosity(); v != silent {
		fmt.Fprintln(os.Stdout, msg)
	}
}

// Warn writes a key message to stderr so as not to interfere with handling json
// on stdout
func Warn(msg string) {
	if v := GetVerbosity(); v != silent && v >= warn {
		fmt.Fprintln(os.Stderr, "[ WARN ] ", msg)
	}
}

// Info writes an extended message to stderr so as not to interfere with handling
// json on stdout
func Info(msg string) {
	if v := GetVerbosity(); v != silent && v >= info {
		fmt.Fprintln(os.Stderr, "[ INFO ] ", msg)
	}
}

// Debug writes a verbose message to stderr so as not to interfere with handling
// json on stdout
func Debug(msg string) {
	if v := GetVerbosity(); v != silent && v >= debug {
		fmt.Fprintln(os.Stderr, "[ DEBUG ]", msg)
	}
}
