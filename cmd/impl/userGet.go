//
// Implementation and supporting functions for the `user get` subcommand.
// Implementation functions make the code more testable
//
package impl

import (
	"encoding/json"
	"site24x7/api"
)

// UserGet is the testable implementation code for cmd.userGetCmd
func UserGet(f UserAccessorFlags, u *api.User, getter func() error) ([]byte, error) {
	if err := f.validate(); err != nil {
		return nil, err
	}

	// Hydrate the user with known values
	u.Id = f.ID
	u.EmailAddress = f.EmailAddress

	if err := getter(); err != nil {
		return nil, err
	}

	out, _ := json.MarshalIndent(u, "", "    ")

	return out, nil
}
