//
// Implementation and supporting functions for the `user get` subcommand.
// Implementation functions make the code more testable
//
package impl

import (
	"fmt"
	"site24x7/api"
	"strings"
	"testing"
)

func TestUserAccessorFlags_validate(t *testing.T) {
	type flags struct {
		Id           string
		EmailAddress string
	}
	tests := []struct {
		name       string
		flags      flags
		wantErr    bool
		wantErrMsg string
	}{
		{
			name: "Errors when all flags are empty",
			flags: flags{
				Id:           "",
				EmailAddress: "",
			},
			wantErr:    true,
			wantErrMsg: "either an ID or an email address is required",
		},
		{
			name: "Only an id value was passed",
			flags: flags{
				Id:           "1001001SOS",
				EmailAddress: "",
			},
			wantErr: false,
		},
		{
			name: "Only an email address value was passed",
			flags: flags{
				Id:           "",
				EmailAddress: "fred@flintstone.com",
			},
			wantErr: false,
		},
		{
			name: "Both an id and an email address were passed",
			flags: flags{
				Id:           "1001001SOS",
				EmailAddress: "fred@flintstone.com",
			},
			wantErr:    true,
			wantErrMsg: "please include either an ID OR",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := UserAccessorFlags{
				ID:           tt.flags.Id,
				EmailAddress: tt.flags.EmailAddress,
			}
			err := f.validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("UserAccessorFlags.validate() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil && !strings.HasPrefix(err.Error(), tt.wantErrMsg) {
				t.Errorf("UserAccessorFlags.validate() error = %v, wantErr %v", err.Error(), tt.wantErrMsg)
			}
		})
	}
}

func TestUserDelete(t *testing.T) {
	type args struct {
		f       UserAccessorFlags
		u       *api.User
		deleter func() error
	}
	tests := []struct {
		name       string
		args       args
		wantErr    bool
		wantErrMsg string
	}{
		{
			name: "Rethrows a flag validation error",
			args: args{
				f: UserAccessorFlags{ID: "", EmailAddress: ""},
				u: &api.User{},
				deleter: func() error {
					return nil
				},
			},
			wantErr:    true,
			wantErrMsg: "either an ID or an email address is required",
		},
		{
			name: "Rethrows a user deleter error",
			args: args{
				f: UserAccessorFlags{ID: "1001001SOS"},
				u: &api.User{},
				deleter: func() error {
					return fmt.Errorf("Whoops!")
				},
			},
			wantErr:    true,
			wantErrMsg: "Whoops!",
		},
		{
			name: "Successfully deletes the user",
			args: args{
				f: UserAccessorFlags{ID: "1001001SOS"},
				u: &api.User{},
				deleter: func() error {
					return nil
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := UserDelete(tt.args.f, tt.args.u, tt.args.deleter)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserDelete() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil && !strings.HasPrefix(err.Error(), tt.wantErrMsg) {
				t.Errorf("UserDelete() error msg = %v, wantErrMsg %v", err.Error(), tt.wantErrMsg)
			}
		})
	}
}
