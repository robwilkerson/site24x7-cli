package logger

import (
	"os"
	"testing"

	"github.com/spf13/pflag"
)

func getDefaultFlags() *pflag.FlagSet {
	fs := pflag.NewFlagSet("verbosityTestFlags", pflag.ExitOnError)
	fs.BoolP("quiet", "q", false, "")
	fs.CountP("verbose", "v", "")

	return fs
}

func TestSetVerbosity(t *testing.T) {
	type args struct {
		fs *pflag.FlagSet
	}

	tests := []struct {
		name  string
		args  args
		flags []string
		want  int
	}{
		{
			name: "Sets verbosity to silent",
			args: args{
				fs: getDefaultFlags(),
			},
			flags: []string{"-q"},
			want:  -1,
		},
		{
			name: "Sets verbosity to normal",
			args: args{
				fs: getDefaultFlags(),
			},
			flags: nil,
			want:  0,
		},
		{
			name: "Sets verbosity to warn",
			args: args{
				fs: getDefaultFlags(),
			},
			flags: []string{"-v"},
			want:  1,
		},
		{
			name: "Sets verbosity to info",
			args: args{
				fs: getDefaultFlags(),
			},
			flags: []string{"-vv"},
			want:  2,
		},
		{
			name: "Sets verbosity to debug",
			args: args{
				fs: getDefaultFlags(),
			},
			flags: []string{"-vvv"},
			want:  3,
		},
		{
			name: "Ignores verbose flags when quiet is set",
			args: args{
				fs: getDefaultFlags(),
			},
			flags: []string{"-vvv", "-q"},
			want:  -1,
		},
	}
	for _, tt := range tests {
		tt.args.fs.Parse(tt.flags)
		t.Run(tt.name, func(t *testing.T) {
			if got := SetVerbosity(tt.args.fs); got != tt.want {
				t.Errorf("SetVerbosity() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetVerbosity(t *testing.T) {
	tests := []struct {
		name   string
		before func()
		want   int
	}{
		{
			name: "Returns quiet mode",
			before: func() {
				os.Setenv("VERBOSITY", "-1")
			},
			want: -1,
		},
		{
			name: "Returns debug",
			before: func() {
				os.Setenv("VERBOSITY", "3")
			},
			want: 3,
		},
		{
			name: "Returns normal in an error case",
			before: func() {
				// noop
				// No env var exists, so the read will error
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		// Reset the relevant env var
		os.Unsetenv("VERBOSITY")

		tt.before()
		t.Run(tt.name, func(t *testing.T) {
			if got := GetVerbosity(); got != tt.want {
				t.Errorf("GetVerbosity() = %v, want %v", got, tt.want)
			}
		})
	}
}
