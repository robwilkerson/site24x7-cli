package user

import (
	"encoding/json"
	"fmt"
	"site24x7/api"
	"site24x7/cmd/impl"

	"strings"

	"github.com/spf13/pflag"
)

// Alias upstream functions for mocking

var apiUserList = api.UserList
var apiUserGet = api.UserGet
var apiUserCreate = api.UserCreate
var apiUserUpdate = api.UserUpdate
var apiUserDelete = api.UserDelete

// list returns a slice containing all users on the account
var list = func() ([]api.User, error) {
	data, err := apiUserList()
	if err != nil {
		return nil, err
	}

	var users []api.User
	if err = json.Unmarshal(data, &users); err != nil {
		return nil, fmt.Errorf("[user.findByEmail] Unable to  parse response data (%s)", err)
	}

	return users, nil
}

// findByEmail returns a user with a given email address
var findByEmail = func(email string) (*api.User, error) {
	users, err := list()
	if err != nil {
		return nil, err
	}

	// Extract the one with a matching email address
	for _, u := range users {
		if strings.EqualFold(u.EmailAddress, email) {
			return &u, nil
		}
	}

	return nil, &api.NotFoundError{Message: fmt.Sprintf("[user.findByEmail] User (%s) not found", email)}
}

// get fetches a user either by email address or by identifier
var get = func(id string, email string) (*api.User, error) {
	var u api.User

	if email != "" {
		// Fetch by email address
		r, err := findByEmail(email)
		if err != nil {
			return nil, err
		}

		// Dereference the returned pointer into our user var
		u = *r
	} else {
		// Fetch by user ID - a.k.a, the official way
		data, err := apiUserGet(id)
		if err != nil {
			return nil, err
		}

		// Ensure that we have a fully hydrated user struct
		if err = json.Unmarshal(data, &u); err != nil {
			return nil, fmt.Errorf("[user.get] Unable to  parse response data (%s)", err)
		}
	}

	return &u, nil
}

// Create is the implementation of the `user create` command
func Create(email string, fs *pflag.FlagSet) ([]byte, error) {
	// Panics if a flag doesn't validate
	validateWriters(fs)

	u := &api.User{EmailAddress: email}
	fs.VisitAll(func(f *pflag.Flag) {
		// StatusIQRole & CloudspendRole may not exist for some accounts and the
		// default value is invalid to ensure that it returns an error. For
		// these we want to explicitly exclude them if they weren't changed.
		if (f.Name == "statusiq-role" || f.Name == "cloudspend-role") && !f.Changed {
			return
		}

		property := normalizeName(f)
		value := impl.TypedFlagValue(fs, f)
		if strings.HasPrefix(property, "AlertingPeriod") {
			impl.SetProperty(&u.AlertSettings.AlertingPeriod, strings.Replace(property, "AlertingPeriod", "", -1), value)
		} else if strings.HasPrefix(property, "Alert") {
			impl.SetProperty(&u.AlertSettings, strings.Replace(property, "Alert", "", -1), value)
		} else if strings.HasPrefix(property, "Mobile") {
			impl.SetProperty(&u.MobileSettings, strings.Replace(property, "Mobile", "", -1), value)
		} else {
			impl.SetProperty(u, property, value)
		}
	})

	data, err := apiUserCreate(u)
	if err != nil {
		return nil, err
	}

	// Ensure that we have a fully hydrated user struct
	var usr api.User
	if err = json.Unmarshal(data, &usr); err != nil {
		return nil, fmt.Errorf("[user.Create] Unable to  parse response data (%s)", err)
	}

	// Return json for display purposes
	j, _ := json.MarshalIndent(usr, "", "    ")

	return j, nil
}

// Get is the implementation of the `user get` command
func Get(fs *pflag.FlagSet) ([]byte, error) {
	validateAccessors(fs)

	id, _ := fs.GetString("id")
	email, _ := fs.GetString("email")

	u, err := get(id, email)
	if err != nil {
		return nil, err
	}

	j, _ := json.MarshalIndent(u, "", "    ")

	return j, nil
}

// Update is the implementation of the `user update` command
func Update(fs *pflag.FlagSet) ([]byte, error) {
	validateAccessors(fs)
	validateWriters(fs)

	id, _ := fs.GetString("id")
	email, _ := fs.GetString("email")
	u, err := get(id, email)
	if err != nil {
		return nil, err
	}

	// Hydrate the user, updating ONLY flags that were set
	fs.Visit(func(f *pflag.Flag) {
		property := normalizeName(f)
		value := impl.TypedFlagValue(fs, f)
		if strings.HasPrefix(property, "AlertingPeriod") {
			impl.SetProperty(&u.AlertSettings.AlertingPeriod, strings.Replace(property, "AlertingPeriod", "", -1), value)
		} else if strings.HasPrefix(property, "Alert") {
			impl.SetProperty(&u.AlertSettings, strings.Replace(property, "Alert", "", -1), value)
		} else if strings.HasPrefix(property, "Mobile") {
			impl.SetProperty(&u.MobileSettings, strings.Replace(property, "Mobile", "", -1), value)
		} else {
			impl.SetProperty(u, property, value)
		}
	})

	data, err := apiUserUpdate(u)
	if err != nil {
		return nil, err
	}

	// Ensure that we have a fully hydrated user struct
	var usr api.User
	if err = json.Unmarshal(data, &usr); err != nil {
		return nil, fmt.Errorf("[user.Update] Unable to  parse response data (%s)", err)
	}

	j, _ := json.MarshalIndent(usr, "", "    ")

	return j, nil
}

// Delete is the implementation of the `user delete` command
func Delete(fs *pflag.FlagSet) error {
	validateAccessors(fs)

	id, _ := fs.GetString("id")
	email, _ := fs.GetString("email")

	u, err := get(id, email)
	if err != nil {
		return err
	}

	if err := apiUserDelete(u.ID); err != nil {
		return err
	}

	return nil
}

// List is the implementation of the `user list` command
func List() ([]byte, error) {
	users, err := list()
	if err != nil {
		return nil, err
	}

	j, _ := json.MarshalIndent(users, "", "    ")

	return j, nil
}
