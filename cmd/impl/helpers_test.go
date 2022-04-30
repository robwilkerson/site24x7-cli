package impl

import (
	"reflect"
	"testing"
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
