package impl

import (
	"encoding/json"
	"fmt"
	"reflect"
	"site24x7/api"
	"testing"
)

func Test_userGetFlags_validate(t *testing.T) {
	type flags struct {
		id           string
		emailAddress string
	}
	tests := []struct {
		name    string
		flags   flags
		wantErr bool
	}{
		{
			name: "All flag values are empty",
			flags: flags{
				id:           "",
				emailAddress: "",
			},
			wantErr: true,
		},
		{
			name: "Only an id value was passed",
			flags: flags{
				id:           "1001001SOS",
				emailAddress: "",
			},
			wantErr: false,
		},
		{
			name: "Only an email address value was passed",
			flags: flags{
				id:           "",
				emailAddress: "fred@flintstone.com",
			},
			wantErr: false,
		},
		{
			name: "Both an id and an email address were passed",
			flags: flags{
				id:           "1001001SOS",
				emailAddress: "fred@flintstone.com",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &UserAccessorFlags{
				ID:           tt.flags.id,
				EmailAddress: tt.flags.emailAddress,
			}
			if err := f.validate(); (err != nil) != tt.wantErr {
				t.Errorf("userGetFlags.validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_userGet(t *testing.T) {
	type args struct {
		f      UserAccessorFlags
		u      *api.User
		getter func() error
	}

	mockUserIn := &api.User{Id: "1001001SOS"}
	mockUserOut, _ := json.MarshalIndent(mockUserIn, "", "    ")
	tests := []struct {
		name       string
		args       args
		want       []byte
		wantErr    bool
		wantErrMsg string
	}{
		{
			name: "Rethrows an invalid flag error",
			args: args{
				f: UserAccessorFlags{ID: "", EmailAddress: ""},
				u: &api.User{},
				getter: func() error {
					return nil
				},
			},
			want:       nil,
			wantErr:    true,
			wantErrMsg: "either an ID or an email address is required to retrieve a user",
		},
		{
			name: "Handles an error thrown by the getter",
			args: args{
				f: UserAccessorFlags{ID: "1001001SOS", EmailAddress: ""},
				u: &api.User{},
				getter: func() error {
					return fmt.Errorf("Whoops!")
				},
			},
			want:       nil,
			wantErr:    true,
			wantErrMsg: "Whoops!",
		},
		{
			name: "Returns no error",
			args: args{
				f: UserAccessorFlags{ID: "1001001SOS", EmailAddress: ""},
				u: mockUserIn,
				getter: func() error {
					return nil
				},
			},
			want:    mockUserOut,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := UserGet(tt.args.f, tt.args.u, tt.args.getter)
			if (err != nil) != tt.wantErr {
				t.Errorf("userGet() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && err.Error() != tt.wantErrMsg {
				t.Errorf("userGet() error msg = %s, wantErrMsg = %s", err.Error(), tt.wantErrMsg)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("userGet() = %v, want %v", got, tt.want)
			}
		})
	}
}
