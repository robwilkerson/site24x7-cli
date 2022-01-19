package user

import (
	"encoding/json"
	"site24x7/api"

	"github.com/spf13/pflag"
)

// Read is the implementation of the `user get` command
func Read(fs *pflag.FlagSet, u *api.User, getter func() error) ([]byte, error) {
	validateAccessors(fs)

	// Hydrate the user with known values
	u.Id, _ = fs.GetString("id")
	u.EmailAddress, _ = fs.GetString("email")

	if err := getter(); err != nil {
		return nil, err
	}

	out, _ := json.MarshalIndent(u, "", "    ")

	return out, nil
}
