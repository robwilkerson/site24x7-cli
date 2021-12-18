/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"reflect"
	"testing"
)

func Test_lookupIds(t *testing.T) {
	type args struct {
		list   []string
		lookup map[string]int
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			name: "Key exists",
			args: args{
				[]string{"b"},
				map[string]int{
					"A": 0,
					"B": 1,
				},
			},
			want: []int{1},
		},
		{
			name: "Key doesn't exist",
			args: args{
				[]string{"d"},
				map[string]int{"A": 0, "B": 1, "C": 2},
			},
			want: nil,
		},
		{
			name: "Multiple keys exist",
			args: args{
				[]string{"b", "d", "y"},
				map[string]int{
					"A": 0,
					"B": 1,
					"c": 2,
					"D": 3,
					"X": 10,
					"Y": 11,
				},
			},
			want: []int{1, 3, 11},
		},
		{
			name: "No keys exist",
			args: args{
				[]string{"m", "l", "p", "r", "t"},
				map[string]int{
					"A": 0,
					"B": 1,
					"C": 2,
					"D": 3,
					"X": 10,
					"Y": 11,
				},
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := lookupIds(tt.args.list, tt.args.lookup); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("lookupIds() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_validateUserGetInput(t *testing.T) {
	type args struct {
		id    string
		email string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Only empty values passed",
			args: args{
				id:    "",
				email: "",
			},
			wantErr: true,
		},
		{
			name: "ID passed",
			args: args{
				id:    "1001001SOS",
				email: "",
			},
			wantErr: false,
		},
		{
			name: "Email passed",
			args: args{
				id:    "",
				email: "fred@flintstone.com",
			},
			wantErr: false,
		},
		{
			name: "ID and Email passed",
			args: args{
				id:    "1001001SOS",
				email: "fred@flintstone.com",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateUserGetInput(tt.args.id, tt.args.email); (err != nil) != tt.wantErr {
				t.Errorf("validateUserGetInput() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
