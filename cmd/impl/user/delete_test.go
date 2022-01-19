package user

import (
	"errors"
	"site24x7/api"
	"testing"

	"github.com/spf13/pflag"
)

func TestDelete(t *testing.T) {
	type args struct {
		fs      *pflag.FlagSet
		u       *api.User
		deleter func() error
	}

	fs := pflag.NewFlagSet("defaultTestFlags", pflag.ExitOnError)
	fs.String("id", "1001001SOS", "")
	fs.String("email", "", "")

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Handles an error thrown by the deleter",
			args: args{
				fs: fs,
				u:  &api.User{},
				deleter: func() error {
					return errors.New("testing!")
				},
			},
			wantErr: true,
		},
		{
			name: "User is deleted successfully",
			args: args{
				fs: fs,
				u:  &api.User{},
				deleter: func() error {
					return nil
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Delete(tt.args.fs, tt.args.u, tt.args.deleter); (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
