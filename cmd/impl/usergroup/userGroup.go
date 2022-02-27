package usergroup

import (
	"encoding/json"
	"fmt"
	"site24x7/api"
)

// Alias upstream functions for mocking

// var apiUserGroupCreate = api.UserGroupCreate
var apiUserGroupGet = api.UserGroupGet

// var apiUserGroupUpdate = api.UserGroupUpdate
// var apiUserGroupDelete = api.UserGroupDelete
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

// Get is the implementation of the `user_group get` command
func Get(id string) ([]byte, error) {
	ug, err := get(id)
	if err != nil {
		return nil, err
	}

	j, _ := json.MarshalIndent(ug, "", "    ")

	return j, nil
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
