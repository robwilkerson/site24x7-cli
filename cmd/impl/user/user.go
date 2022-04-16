package user

import (
	"encoding/json"
	"fmt"
	"reflect"
	"site24x7/api"
	"site24x7/logger"
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

// setProperty sets either a user property or a property on one of a user's
// nested property structures.
func setProperty(v any, property string, value any) {
	logger.Debug(fmt.Sprintf("Setting %s; value: %v\n", property, value))

	rv := reflect.ValueOf(v)

	// dereference the pointer
	rv = rv.Elem()

	// lookup the field by name and set the new value
	f := rv.FieldByName(property)

	if f.IsValid() {
		f.Set(reflect.ValueOf(value))
	} else {
		logger.Warn(fmt.Sprintf("[usergroup.setProperty] Invalid user group property %s; ignoring", property))
	}
}

// Create is the implementation of the `user create` command
func Create(email string, fs *pflag.FlagSet) ([]byte, error) {
	// Panics if a flag doesn't validate
	validateWriters(fs)

	u := &api.User{EmailAddress: email}
	fs.VisitAll(func(f *pflag.Flag) {
		// If this is a flag that doesn't directly map to a user property,
		// skip it by returning early
		if _, ok := nonUserFlags[f.Name]; ok {
			return
		}

		// StatusIQRole & CloudspendRole may not exist for some accounts and the
		// default value is invalid to ensure that it returns an error. For
		// these we want to explicitly exclude them if they weren't changed.
		if (f.Name == "statusiq-role" || f.Name == "cloudspend-role") && !f.Changed {
			return
		}

		// Extract the appropriately typed value from the flag
		var v any
		switch f.Value.Type() {
		case "string":
			v, _ = fs.GetString(f.Name)
		case "int":
			v, _ = fs.GetInt(f.Name)
		case "stringSlice":
			v, _ = fs.GetStringSlice(f.Name)
		case "intSlice":
			v, _ = fs.GetIntSlice(f.Name)
		case "bool":
			v, _ = fs.GetBool(f.Name)
		default:
			// This is a problem, but I'm not sure it needs to be a fatal one
			logger.Warn(fmt.Sprintf("[user.Create] Unhandled data type (%s) for the %s flag", f.Value.Type(), f.Name))
		}

		// normalize property name
		p := normalizeName(f)

		if strings.HasPrefix(p, "AlertingPeriod") {
			setProperty(&u.AlertSettings.AlertingPeriod, strings.Replace(p, "AlertingPeriod", "", -1), v)
		} else if strings.HasPrefix(p, "Alert") {
			setProperty(&u.AlertSettings, strings.Replace(p, "Alert", "", -1), v)
		} else if strings.HasPrefix(p, "Mobile") {
			setProperty(&u.MobileSettings, strings.Replace(p, "Mobile", "", -1), v)
		} else {
			setProperty(u, p, v)
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
		// If this is a flag that doesn't directly map to a user property,
		// skip it by returning early
		if _, ok := nonUserFlags[f.Name]; ok {
			return
		}

		// Extract the appropriately typed value from the flag
		var v any
		switch f.Value.Type() {
		case "string":
			v, _ = fs.GetString(f.Name)
		case "int":
			v, _ = fs.GetInt(f.Name)
		case "stringSlice":
			v, _ = fs.GetStringSlice(f.Name)
		case "intSlice":
			v, _ = fs.GetIntSlice(f.Name)
		case "bool":
			v, _ = fs.GetBool(f.Name)
		default:
			// This is a problem, but I'm not sure it needs to be a fatal one
			logger.Warn(fmt.Sprintf("[user.Update] Unhandled data type (%s) for the %s flag", f.Value.Type(), f.Name))
		}

		// normalize property name
		p := normalizeName(f)

		if strings.HasPrefix(p, "AlertingPeriod") {
			setProperty(&u.AlertSettings.AlertingPeriod, strings.Replace(p, "AlertingPeriod", "", -1), v)
		} else if strings.HasPrefix(p, "Alert") {
			setProperty(&u.AlertSettings, strings.Replace(p, "Alert", "", -1), v)
		} else if strings.HasPrefix(p, "Mobile") {
			setProperty(&u.MobileSettings, strings.Replace(p, "Mobile", "", -1), v)
		} else {
			setProperty(u, p, v)
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
