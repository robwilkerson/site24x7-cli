//
// Implementation and supporting functions for the `user get` subcommand.
// Implementation functions make the code more testable
//
package impl

import (
	"site24x7/api"
)

// UserDelete is the testable implementation code for cmd.userDeleteCmd
func UserDelete(f UserAccessorFlags, u *api.User, deleter func() error) error {
	if err := f.validate(); err != nil {
		return err
	}

	// Hydrate the user with known values
	u.Id = f.ID
	u.EmailAddress = f.EmailAddress

	if err := deleter(); err != nil {
		return err
	}

	return nil
}
