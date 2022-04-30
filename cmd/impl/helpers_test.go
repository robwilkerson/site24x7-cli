package impl

import (
	"reflect"
	"testing"

	"github.com/spf13/pflag"
)

func TestSetProperty(t *testing.T) {
	type args struct {
		v        any
		property string
		value    any
	}
	type mockStruct struct {
		MockString string
		MockInt    int
		private    string
	}
	original := mockStruct{"Studio", 54, "Benjamin"}
	expected := mockStruct{"Area", 54, "Benjamin"}
	tests := []struct {
		name     string
		args     args
		want     bool
		expected mockStruct
	}{
		{
			name: "Ignores any input that's not a pointer to a struct",
			args: args{
				v:        "Not a struct",
				property: "foo",
				value:    "bar",
			},
			want:     false,
			expected: original,
		},
		{
			name: "Ignores a property that doesn't exist",
			args: args{
				v:        &original,
				property: "someOtherString",
				value:    "ignoreme",
			},
			want:     false,
			expected: original,
		},
		{
			name: "Does not set a property that isn't exported",
			args: args{
				v:        &original,
				property: "private",
				value:    "Pyle",
			},
			want:     false,
			expected: original,
		},
		{
			name: "Updates a valid struct property",
			args: args{
				v:        &original,
				property: "MockString",
				value:    "Area",
			},
			want:     true,
			expected: expected,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SetProperty(tt.args.v, tt.args.property, tt.args.value); got != tt.want {
				t.Errorf("SetProperty() = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(original, tt.expected) {
				t.Errorf("SetProperty() got = %v, expected %v", original, tt.expected)
			}
		})
	}
}

func TestTypedFlagValue(t *testing.T) {
	type args struct {
		fs *pflag.FlagSet
		f  *pflag.Flag
	}
	mockFlagSet := pflag.NewFlagSet("Testing", pflag.PanicOnError)
	mockFlagSet.String("string", "string value", "")
	mockFlagSet.StringSlice("stringslice", []string{"Slice", "of", "Strings"}, "")
	mockFlagSet.Int("int", 19, "")
	mockFlagSet.IntSlice("intslice", []int{1, 2, 3, 5, 8}, "")
	mockFlagSet.Bool("bool", false, "")
	mockFlagSet.Float64("unhandled", 2.3980959, "")
	tests := []struct {
		name string
		args args
		want any
	}{
		{
			name: "Returns a string value",
			args: args{
				fs: mockFlagSet,
				f:  mockFlagSet.Lookup("string"),
			},
			want: "string value",
		},
		{
			name: "Returns a string slice value",
			args: args{
				fs: mockFlagSet,
				f:  mockFlagSet.Lookup("stringslice"),
			},
			want: []string{"Slice", "of", "Strings"},
		},
		{
			name: "Returns an int value",
			args: args{
				fs: mockFlagSet,
				f:  mockFlagSet.Lookup("int"),
			},
			want: 19,
		},
		{
			name: "Returns a int slice value",
			args: args{
				fs: mockFlagSet,
				f:  mockFlagSet.Lookup("intslice"),
			},
			want: []int{1, 2, 3, 5, 8},
		},
		{
			name: "Returns an bool value",
			args: args{
				fs: mockFlagSet,
				f:  mockFlagSet.Lookup("bool"),
			},
			want: false,
		},
		{
			name: "Returns nothing",
			args: args{
				fs: mockFlagSet,
				f:  mockFlagSet.Lookup("unhandled"),
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := TypedFlagValue(tt.args.fs, tt.args.f); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TypedFlagValue() = %v, want %v", got, tt.want)
			}
		})
	}
}
