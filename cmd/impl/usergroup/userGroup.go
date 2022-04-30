package usergroup

import (
	"encoding/json"
	"fmt"
	"site24x7/api"
	"site24x7/cmd/impl"
	"site24x7/logger"

	"github.com/spf13/pflag"
)

// Alias upstream functions for mocking

var apiUserGroupCreate = api.UserGroupCreate
var apiUserGroupGet = api.UserGroupGet

// var apiUserGroupUpdate = api.UserGroupUpdate
var apiUserGroupDelete = api.UserGroupDelete
var apiUserGroupList = api.UserGroupList

// list returns a slice containing all users on the account
var list = func() ([]api.UserGroup, error) {
	data, err := apiUserGroupList()
	if err != nil {
		return nil, err
	}

	var list []api.UserGroup
	if err = json.Unmarshal(data, &list); err != nil {
		return nil, fmt.Errorf("[usergroup.list] Unable to  parse response data (%s)", err)
	}

	return list, nil
}

// get fetches a user group
var get = func(id string) (*api.UserGroup, error) {
	var ug api.UserGroup

	data, err := apiUserGroupGet(id)
	if err != nil {
		return nil, err
	}

	// Ensure that we have a fully hydrated struct
	if err = json.Unmarshal(data, &ug); err != nil {
		return nil, fmt.Errorf("[usergroup.get] Unable to  parse response data (%s)", err)
	}

	return &ug, nil
}

// Create is the implementation of the `user_group create` command
func Create(name string, fs *pflag.FlagSet) ([]byte, error) {
	ug := &api.UserGroup{Name: name}
	fs.VisitAll(func(f *pflag.Flag) {
		// Extract the appropriately typed value from the flag
		var v any
		switch f.Value.Type() {
		case "string":
			v, _ = fs.GetString(f.Name)
		case "int":
			v, _ = fs.GetInt(f.Name)
		case "stringSlice":
			v, _ = fs.GetStringSlice(f.Name)
		default:
			// This is a problem, but I'm not sure it needs to be a fatal one
			logger.Warn(fmt.Sprintf("[usergroup.Create] Unhandled data type (%s) for the %s flag; ignoring", f.Value.Type(), f.Name))
			return
		}

		// normalize property name
		p := normalizeName(f)

		impl.SetProperty(ug, p, v)
	})

	data, err := apiUserGroupCreate(ug)
	if err != nil {
		return nil, err
	}

	// Ensure that we have a fully hydrated user struct
	var usergru api.UserGroup
	if err = json.Unmarshal(data, &usergru); err != nil {
		return nil, fmt.Errorf("[usergroup.Create] Unable to  parse response data (%s)", err)
	}

	// Return json for display purposes
	j, _ := json.MarshalIndent(usergru, "", "    ")

	return j, nil
}

// Get is the implementation of the `user_group get` command
func Get(id string) ([]byte, error) {
	ug, err := get(id)
	if err != nil {
		return nil, err
	}

	j, _ := json.MarshalIndent(ug, "", "    ")

	return j, nil
}

// Delete is the implementation of the `monitor_group delete` command
func Delete(id string) error {
	err := apiUserGroupDelete(id)
	if err != nil {
		return err
	}

	return nil
}

// List is the implementation of the `user_group list` command
func List() ([]byte, error) {
	list, err := list()
	if err != nil {
		return nil, err
	}

	j, _ := json.MarshalIndent(list, "", "    ")

	return j, nil
}
