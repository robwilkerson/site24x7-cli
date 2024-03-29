package monitorgroup

import (
	"encoding/json"
	"fmt"
	"site24x7/api"
	"site24x7/cmd/impl"
	"site24x7/logger"

	"github.com/spf13/pflag"
)

// Alias upstream functions for mocking

var apiMonitorGroupList = api.MonitorGroupList
var apiMonitorGroupGet = api.MonitorGroupGet
var apiMonitorGroupCreate = api.MonitorGroupCreate
var apiMonitorGroupUpdate = api.MonitorGroupUpdate
var apiMonitorGroupDelete = api.MonitorGroupDelete

// list returns a slice containing all users on the account
var list = func(withSubgroups bool) ([]api.MonitorGroup, error) {
	data, err := apiMonitorGroupList(withSubgroups)
	if err != nil {
		return nil, err
	}

	var mongrus []api.MonitorGroup
	if err = json.Unmarshal(data, &mongrus); err != nil {
		return nil, fmt.Errorf("[monitorgroup.list] Unable to  parse response data (%s)", err)
	}

	return mongrus, nil
}

// get fetches a monitor group
var get = func(id string) (*api.MonitorGroup, error) {
	var mg api.MonitorGroup

	data, err := apiMonitorGroupGet(id)
	if err != nil {
		return nil, err
	}

	// Ensure that we have a fully hydrated struct
	if err = json.Unmarshal(data, &mg); err != nil {
		return nil, fmt.Errorf("[monitorgroup.get] Unable to  parse response data (%s)", err)
	}

	return &mg, nil
}

// Create is the implementation of the `monitor_group create` command
func Create(name string, fs *pflag.FlagSet) ([]byte, error) {
	mg := &api.MonitorGroup{Name: name}
	fs.VisitAll(func(f *pflag.Flag) {
		property := normalizeName(f)
		value := impl.TypedFlagValue(fs, f)

		impl.SetProperty(mg, property, value)
	})

	data, err := apiMonitorGroupCreate(mg)
	if err != nil {
		return nil, err
	}

	// Ensure that we have a fully hydrated user struct
	var mongru api.MonitorGroup
	if err = json.Unmarshal(data, &mongru); err != nil {
		return nil, fmt.Errorf("[monitorGroup.Create] Unable to  parse response data (%s)", err)
	}

	// Return json for display purposes
	j, _ := json.MarshalIndent(mongru, "", "    ")

	return j, nil
}

// Get is the implementation of the `monitor_group get` command
func Get(id string, fs *pflag.FlagSet) ([]byte, error) {
	mg, err := get(id)
	if err != nil {
		return nil, err
	}

	j, _ := json.MarshalIndent(mg, "", "    ")

	return j, nil
}

// Update is the implementation of the `monitor_group update` command
func Update(id string, fs *pflag.FlagSet) ([]byte, error) {
	logger.Info(fmt.Sprintf("[MonitorGroup.Update] Updating group with ID %s", id))

	mg, err := get(id)
	if err != nil {
		return nil, err
	}

	logger.Debug(fmt.Sprintf("[MonitorGroup.Update] Fetched group %+v", mg))

	// Hydrate the user, updating ONLY flags that were set
	fs.Visit(func(f *pflag.Flag) {
		property := normalizeName(f)
		value := impl.TypedFlagValue(fs, f)

		impl.SetProperty(mg, property, value)
	})

	data, err := apiMonitorGroupUpdate(mg)
	if err != nil {
		return nil, err
	}

	// Ensure that we have a fully hydrated user struct
	var mgOut api.MonitorGroup
	if err = json.Unmarshal(data, &mgOut); err != nil {
		return nil, fmt.Errorf("[monitorgroup.Update] Unable to  parse response data (%s)", err)
	}

	j, _ := json.MarshalIndent(mgOut, "", "    ")

	return j, nil
}

// Delete is the implementation of the `monitor_group delete` command
func Delete(id string, fs *pflag.FlagSet) error {
	err := apiMonitorGroupDelete(id)
	if err != nil {
		return err
	}

	return nil
}

// List is the implementation of the `monitor_group list` command
func List(fs *pflag.FlagSet) ([]byte, error) {
	sg, _ := fs.GetBool("with-subgroups")

	mongrus, err := list(sg)
	if err != nil {
		return nil, err
	}

	j, _ := json.MarshalIndent(mongrus, "", "    ")

	return j, nil
}
