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

var apiUserGroupUpdate = api.UserGroupUpdate
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
		property := normalizeName(f)
		value := impl.TypedFlagValue(fs, f)

		impl.SetProperty(ug, property, value)
	})

	data, err := apiUserGroupCreate(ug)
	if err != nil {
		return nil, err
	}

	// Ensure that we have a fully hydrated user group struct
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

// Update is the implementation of the `user_group update` command
func Update(id string, fs *pflag.FlagSet) ([]byte, error) {
	logger.Info(fmt.Sprintf("[UserGroup.Update] Updating group with ID %s", id))

	ug, err := get(id)
	if err != nil {
		return nil, err
	}

	logger.Debug(fmt.Sprintf("[UserGroup.Update] Fetched group %+v", ug))

	// Hydrate the user, updating ONLY flags that were set
	fs.Visit(func(f *pflag.Flag) {
		property := normalizeName(f)
		value := impl.TypedFlagValue(fs, f)

		impl.SetProperty(ug, property, value)
	})

	data, err := apiUserGroupUpdate(ug)
	if err != nil {
		return nil, err
	}

	// Ensure that we have a fully hydrated user struct
	var ugOut api.UserGroup
	if err = json.Unmarshal(data, &ugOut); err != nil {
		return nil, fmt.Errorf("[monitorgroup.Update] Unable to  parse response data (%s)", err)
	}

	j, _ := json.MarshalIndent(ugOut, "", "    ")

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
