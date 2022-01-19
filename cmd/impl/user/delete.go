package user

import (
	"site24x7/api"

	"github.com/spf13/pflag"
)

func Delete(fs *pflag.FlagSet, u *api.User, deleter func() error) error {
	validateAccessors(fs)

	// Hydrate the user with known values
	u.Id, _ = fs.GetString("id")
	u.EmailAddress, _ = fs.GetString("email")

	if err := deleter(); err != nil {
		return err
	}

	return nil
}
