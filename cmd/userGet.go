package cmd

import (
	"encoding/json"
	"fmt"
	"site24x7/api"
)

// flags contains the value of any flag sent to the command
type userGetFlags struct {
	id           string
	emailAddress string
}

// validate validates user data passed to the get command
func (f userGetFlags) validate() error {
	if f.id != "" && f.emailAddress != "" {
		return fmt.Errorf("please include either an ID OR an email address, not both")
	} else if f.id == "" && f.emailAddress == "" {
		return fmt.Errorf("either an ID or an email address is required to retrieve a user")
	}

	return nil
}

// userGet is the testable implementation code for userGetCmd
func userGet(f userGetFlags, u *api.User, getter func() error) ([]byte, error) {
	if err := f.validate(); err != nil {
		return nil, err
	}

	// Hydrate the user with known values
	u.Id = f.id
	u.EmailAddress = f.emailAddress

	if err := getter(); err != nil {
		return nil, err
	}

	out, _ := json.MarshalIndent(u, "", "    ")

	return out, nil
}
