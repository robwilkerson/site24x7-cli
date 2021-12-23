//
// Implementation and supporting functions for the `user get` subcommand.
// Implementation functions make the code more testable
//
package impl

import (
	"fmt"
	"site24x7/api"
)

// UserDeleteFlags contains the value of any flag sent to the command
// TODO: Maybe try sharing userAccessorFlags and userCreatorFlags? These flags are
// exactly the same b/c the serve the same purpose.
type UserDeleteFlags struct {
	Id           string
	EmailAddress string
}

// TODO: Create a type validator interface{} to combine some of these
// duplicative operations

// validate validates user data passed to the `user delete` command
func (f UserDeleteFlags) validate() error {
	if f.Id != "" && f.EmailAddress != "" {
		return fmt.Errorf("please include either an ID OR an email address, not both")
	} else if f.Id == "" && f.EmailAddress == "" {
		return fmt.Errorf("either an ID or an email address is required to retrieve a user")
	}

	return nil
}

// UserDelete is the testable implementation code for cmd.userDeleteCmd
func UserDelete(f UserDeleteFlags, u *api.User, deleter func() error) error {
	if err := f.validate(); err != nil {
		return err
	}

	// Hydrate the user with known values
	u.Id = f.Id
	u.EmailAddress = f.EmailAddress

	if err := deleter(); err != nil {
		return err
	}

	return nil
}
