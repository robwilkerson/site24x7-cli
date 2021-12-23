//
// Implementation and supporting functions for the `user list` subcommand.
// Implementation functions make the code more testable
//
package impl

import (
	"encoding/json"
	"fmt"
	"reflect"
	"site24x7/api"
	"strings"
	"testing"
)

func TestUserList(t *testing.T) {
	type args struct {
		getter func() ([]api.User, error)
	}
	mockUsers := []api.User{
		{EmailAddress: "humpty@dumpty.com"},
		{EmailAddress: "alice@wonderland.com"},
	}
	mockUsersOut, _ := json.MarshalIndent(mockUsers, "", "    ")
	tests := []struct {
		name       string
		args       args
		want       []byte
		wantErr    bool
		wantErrMsg string
	}{
		{
			name: "Handles an error thrown by the getter",
			args: args{
				getter: func() ([]api.User, error) {
					return nil, fmt.Errorf("Whoops!")
				},
			},
			want:       nil,
			wantErr:    true,
			wantErrMsg: "Whoops!",
		},
		{
			name: "Returns a list of users",
			args: args{
				getter: func() ([]api.User, error) {
					return mockUsers, nil
				},
			},
			want:    mockUsersOut,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := UserList(tt.args.getter)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && !strings.HasPrefix(err.Error(), tt.wantErrMsg) {
				t.Errorf("UserList() error msg = %s, wantErrMsg = %s", err.Error(), tt.wantErrMsg)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UserList() = %v, want %v", got, tt.want)
			}
		})
	}
}
