package user

import (
	"encoding/json"
	"site24x7/api"

	"github.com/spf13/pflag"
)

// List is the implementation of the `user list` command
func List(getter func() ([]api.User, error)) ([]byte, error) {
	users, err := getter()
	if err != nil {
		return nil, err
	}

	json, _ := json.MarshalIndent(users, "", "    ")

	return json, nil
}

// Read is the implementation of the `user get` command
func Read(fs *pflag.FlagSet, u *api.User, getter func() error) ([]byte, error) {
	validateAccessors(fs)

	// Hydrate the user with known values
	u.Id, _ = fs.GetString("id")
	u.EmailAddress, _ = fs.GetString("email")

	if err := getter(); err != nil {
		return nil, err
	}

	json, _ := json.MarshalIndent(u, "", "    ")

	return json, nil
}
