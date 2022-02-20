package monitorgroup

import (
	"encoding/json"
	"fmt"
	"reflect"
	"site24x7/api"
	"site24x7/logger"

	"github.com/spf13/pflag"
)

// Alias upstream functions for mocking

var apiMonitorGroupList = api.MonitorGroupList
var apiMonitorGroupGet = api.MonitorGroupGet
var apiMonitorGroupCreate = api.MonitorGroupCreate

// var apiUserUpdate = api.UserUpdate
// var apiUserDelete = api.UserDelete

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

// setProperty sets either a user property or a property on one of a user's
// nested property structures.
func setProperty(v interface{}, property string, value interface{}) {
	logger.Debug(fmt.Sprintf("Setting %s; value: %v", property, value))

	rv := reflect.ValueOf(v)

	// dereference the pointer
	rv = rv.Elem()

	// lookup the field by name and set the new value
	f := rv.FieldByName(property)

	if f.IsValid() {
		f.Set(reflect.ValueOf(value))
	} else {
		logger.Debug(fmt.Sprintf("[monitorGroup.setProperty] Invalid monitor group property %s; ignoring", property))
	}

}

// Create is the implementation of the `user create` command
func Create(name string, fs *pflag.FlagSet) ([]byte, error) {
	mg := &api.MonitorGroup{Name: name}
	fs.VisitAll(func(f *pflag.Flag) {
		// Extract the appropriately typed value from the flag
		var v interface{}
		switch f.Value.Type() {
		case "string":
			v, _ = fs.GetString(f.Name)
		case "int":
			v, _ = fs.GetInt(f.Name)
		case "stringSlice":
			v, _ = fs.GetStringSlice(f.Name)
		case "bool":
			v, _ = fs.GetBool(f.Name)
		default:
			// This is a problem, but I'm not sure it needs to be a fatal one
			logger.Warn(fmt.Sprintf("[monitorGroup.Create] Unhandled data type (%s) for the %s flag; ignoring", f.Value.Type(), f.Name))
			return
		}

		// normalize property name
		p := normalizeName(f)

		setProperty(mg, p, v)
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

// Get is the implementation of the `user get` command
func Get(id string, fs *pflag.FlagSet) ([]byte, error) {
	u, err := apiMonitorGroupGet(id)
	if err != nil {
		return nil, err
	}

	j, _ := json.MarshalIndent(u, "", "    ")

	return j, nil
}

// Update is the implementation of the `user update` command
// func Update(fs *pflag.FlagSet) ([]byte, error) {
// 	validateAccessors(fs)
// 	validateWriters(fs)

// 	id, _ := fs.GetString("id")
// 	email, _ := fs.GetString("email")
// 	u, err := get(id, email)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Hydrate the user, updating ONLY flags that were set
// 	fs.Visit(func(f *pflag.Flag) {
// 		// If this is a flag that doesn't directly map to a user property,
// 		// skip it by returning early
// 		if _, ok := nonUserFlags[f.Name]; ok {
// 			return
// 		}

// 		// Extract the appropriately typed value from the flag
// 		var v interface{}
// 		switch f.Value.Type() {
// 		case "string":
// 			v, _ = fs.GetString(f.Name)
// 		case "int":
// 			v, _ = fs.GetInt(f.Name)
// 		case "stringSlice":
// 			v, _ = fs.GetStringSlice(f.Name)
// 		case "intSlice":
// 			v, _ = fs.GetIntSlice(f.Name)
// 		default:
// 			// This is a problem, but I'm not sure it needs to be a fatal one
// 			logger.Warn(fmt.Sprintf("[user.Update] Unhandled data type (%s) for the %s flag", f.Value.Type(), f.Name))
// 		}

// 		// normalize property name
// 		p := normalizeName(f)

// 		if strings.HasPrefix(p, "AlertingPeriod") {
// 			setProperty(&u.AlertSettings.AlertingPeriod, strings.Replace(p, "AlertingPeriod", "", -1), v)
// 		} else if strings.HasPrefix(p, "Alert") {
// 			setProperty(&u.AlertSettings, strings.Replace(p, "Alert", "", -1), v)
// 		} else if strings.HasPrefix(p, "Mobile") {
// 			setProperty(&u.MobileSettings, strings.Replace(p, "Mobile", "", -1), v)
// 		} else {
// 			setProperty(u, p, v)
// 		}
// 	})

// 	data, err := apiUserUpdate(u)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Ensure that we have a fully hydrated user struct
// 	var usr api.User
// 	if err = json.Unmarshal(data, &usr); err != nil {
// 		return nil, fmt.Errorf("[user.Update] Unable to  parse response data (%s)", err)
// 	}

// 	j, _ := json.MarshalIndent(usr, "", "    ")

// 	return j, nil
// }

// Delete is the implementation of the `user delete` command
// func Delete(fs *pflag.FlagSet) error {
// 	validateAccessors(fs)

// 	id, _ := fs.GetString("id")
// 	email, _ := fs.GetString("email")

// 	u, err := get(id, email)
// 	if err != nil {
// 		return err
// 	}

// 	if err := apiUserDelete(u.ID); err != nil {
// 		return err
// 	}

// 	return nil
// }

// List is the implementation of the `user list` command
func List(fs *pflag.FlagSet) ([]byte, error) {
	sg, _ := fs.GetBool("with-subgroups")

	mongrus, err := list(sg)
	if err != nil {
		return nil, err
	}

	j, _ := json.MarshalIndent(mongrus, "", "    ")

	return j, nil
}
