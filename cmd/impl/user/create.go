package user

import (
	"encoding/json"
	"fmt"
	"reflect"
	"site24x7/api"
	"strings"

	"github.com/spf13/pflag"
)

// setProperty sets either a user property or a property on one of a user's
// nested property structures.
func setProperty(v interface{}, property string, value interface{}) {
	// fmt.Printf("Setting %s; value: %v\n", property, value)
	rv := reflect.ValueOf(v)

	// dereference the pointer
	rv = rv.Elem()

	// lookup the field by name and set the new value
	f := rv.FieldByName(property)
	f.Set(reflect.ValueOf(value))
}

// Create is the implementation of the `user create` command
func Create(fs *pflag.FlagSet, u *api.User, creator func() error) ([]byte, error) {
	// Panics if a flag doesn't validate
	validateWriters(fs)

	// Hydrate the user
	ap := &api.UserAlertingPeriod{}
	as := &api.UserAlertSettings{}
	ms := &api.UserMobileSettings{}
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
		var v interface{}
		switch f.Value.Type() {
		case "string":
			v, _ = fs.GetString(f.Name)
		case "int":
			v, _ = fs.GetInt(f.Name)
		case "stringSlice":
			v, _ = fs.GetStringSlice(f.Name)
		case "intSlice":
			v, _ = fs.GetIntSlice(f.Name)
		default:
			// we can't return an error from this function, which would be
			// nice and tidy, so just panic; we def don't want to continue
			panic(fmt.Sprintf("[Create] Unhandled data type (%s) for the %s flag", f.Value.Type(), f.Name))
		}

		// normalize property name
		p := normalizeName(f)

		if strings.HasPrefix(p, "AlertingPeriod") {
			setProperty(ap, strings.Replace(p, "AlertingPeriod", "", -1), v)
		} else if strings.HasPrefix(p, "Alert") {
			setProperty(as, strings.Replace(p, "Alert", "", -1), v)
		} else if strings.HasPrefix(p, "Mobile") {
			setProperty(ms, strings.Replace(p, "Mobile", "", -1), v)
		} else {
			setProperty(u, p, v)
		}
	})

	// Assemble the full user struct
	u.AlertSettings = *as
	u.AlertSettings.AlertingPeriod = *ap
	u.MobileSettings = *ms

	if err := creator(); err != nil {
		return nil, err
	}

	// Return json for display purposes
	json, _ := json.MarshalIndent(u, "", "    ")

	return json, nil
}
