//
// Implementation and supporting functions for the `user get` subcommand.
// Implementation functions make the code more testable
//
package impl

import (
	"encoding/json"
	"fmt"
	"site24x7/api"
)

// userGetFlags contains the value of any flag sent to the command
type UserGetFlags struct {
	Id           string
	EmailAddress string
}

// validate validates user data passed to the get command
func (f UserGetFlags) validate() error {
	if f.Id != "" && f.EmailAddress != "" {
		return fmt.Errorf("please include either an ID OR an email address, not both")
	} else if f.Id == "" && f.EmailAddress == "" {
		return fmt.Errorf("either an ID or an email address is required to retrieve a user")
	}

	return nil
}

// userGet is the testable implementation code for userGetCmd
func UserGet(f UserGetFlags, u *api.User, getter func() error) ([]byte, error) {
	if err := f.validate(); err != nil {
		return nil, err
	}

	// Hydrate the user with known values
	u.Id = f.Id
	u.EmailAddress = f.EmailAddress

	if err := getter(); err != nil {
		return nil, err
	}

	out, _ := json.MarshalIndent(u, "", "    ")

	return out, nil
}
