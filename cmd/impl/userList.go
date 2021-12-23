//
// Implementation and supporting functions for the `user list` subcommand.
// Implementation functions make the code more testable
//
package impl

import (
	"encoding/json"
	"site24x7/api"
)

// userList is the testable implementation code for cmd.userListCmd
func UserList(getter func() ([]api.User, error)) ([]byte, error) {
	users, err := getter()
	if err != nil {
		return nil, err
	}

	out, _ := json.MarshalIndent(users, "", "    ")

	return out, nil
}
