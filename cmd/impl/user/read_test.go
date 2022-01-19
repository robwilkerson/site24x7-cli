package user

import (
	"encoding/json"
	"errors"
	"reflect"
	"site24x7/api"
	"testing"

	"github.com/spf13/pflag"
)

func TestRead(t *testing.T) {
	type args struct {
		fs     *pflag.FlagSet
		u      *api.User
		getter func() error
	}

	d := pflag.NewFlagSet("defaultTestFlags", pflag.ExitOnError)
	d.String("id", "1001001SOS", "")
	d.String("email", "", "")
	userIn := &api.User{}
	userOut := userIn
	userOut.Id = "1001001SOS"
	userOutJSON, _ := json.MarshalIndent(userOut, "", "    ")
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "Handles an error thrown by the getter",
			args: args{
				fs: &pflag.FlagSet{},
				u:  userIn,
				getter: func() error {
					return errors.New("testing!")
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Returns a user",
			args: args{
				fs: d,
				u:  userIn,
				getter: func() error {
					return nil
				},
			},
			want:    userOutJSON,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Read(tt.args.fs, tt.args.u, tt.args.getter)
			if (err != nil) != tt.wantErr {
				t.Errorf("Read() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Read() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestList(t *testing.T) {
	type args struct {
		getter func() ([]api.User, error)
	}
	mockUsers := []api.User{
		{EmailAddress: "humpty@dumpty.com"},
		{EmailAddress: "alice@wonderland.com"},
	}
	mockUsersJSON, _ := json.MarshalIndent(mockUsers, "", "    ")
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "Handles an error thrown by the getter",
			args: args{
				getter: func() ([]api.User, error) {
					return nil, errors.New("testing!")
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Returns a list of users",
			args: args{
				getter: func() ([]api.User, error) {
					return mockUsers, nil
				},
			},
			want:    mockUsersJSON,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := List(tt.args.getter)
			if (err != nil) != tt.wantErr {
				t.Errorf("List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("List() = %v, want %v", got, tt.want)
			}
		})
	}
}
