//
// Implementation and supporting functions for the `user update` subcommand.
// Implementation functions make the code more testable
//
package impl

import (
	"fmt"
	"reflect"
	"site24x7/api"
	"strings"

	"github.com/spf13/pflag"
)

// TODO: These set*Property() functions are effectively duplicates

// setAlertingPeriodProperty sets a specific struct property
func setAlertingPeriodProperty(u *api.UserAlertingPeriod, property string, value interface{}) {
	ru := reflect.ValueOf(u)

	// dereference the pointer
	ru = ru.Elem()

	// lookup the field by name and set the new value
	f := ru.FieldByName(property)
	f.Set(reflect.ValueOf(value))
}

// setAlertSettingsProperty sets a specific struct property
func setAlertSettingsProperty(u *api.UserAlertSettings, property string, value interface{}) {
	ru := reflect.ValueOf(u)

	// dereference the pointer
	ru = ru.Elem()

	// lookup the field by name and set the new value
	f := ru.FieldByName(property)
	f.Set(reflect.ValueOf(value))
}

// setMobileSettingsProperty sets a specific struct property
func setMobileSettingsProperty(u *api.UserMobileSettings, property string, value interface{}) {
	ru := reflect.ValueOf(u)

	// dereference the pointer
	ru = ru.Elem()

	// lookup the field by name and set the new value
	f := ru.FieldByName(property)
	f.Set(reflect.ValueOf(value))
}

// setUserProperty sets a specific struct property
func setUserProperty(u *api.User, property string, value interface{}) {
	ru := reflect.ValueOf(u)

	// dereference the pointer
	ru = ru.Elem()

	// lookup the field by name and set the new value
	f := ru.FieldByName(property)
	f.Set(reflect.ValueOf(value))
}

// UserUpdate is the testable implementation code for cmd.userUpdateCmd
func UserUpdate(f *pflag.FlagSet, u *api.User, updater func() error) error {
	a := UserAccessorFlags{}
	a.ID, _ = f.GetString("ID")
	a.EmailAddress, _ = f.GetString("Email")
	if err := a.validate(); err != nil {
		return err
	}

	// Initialize and fetch the existing user details
	if a.EmailAddress != "" {
		u.EmailAddress = a.EmailAddress
	} else {
		u.Id = a.ID
	}
	if err := u.Get(); err != nil {
		return err
	}

	// Iterate over flags that were explicitly set and update the appropriate
	// user property
	f.Visit(func(fl *pflag.Flag) {
		// Ignore the accessor flags we've extracted above; they're read only
		if fl.Name != "ID" && fl.Name != "Email" {
			// Extract the appropriately typed value from the flag
			var v interface{}
			switch fl.Value.Type() {
			case "string":
				v, _ = f.GetString(fl.Name)
			case "int":
				v, _ = f.GetInt(fl.Name)
			case "stringSlice":
				v, _ = f.GetStringSlice(fl.Name)
			case "intSlice":
				v, _ = f.GetIntSlice(fl.Name)
			default:
				// we can't return an error from this function, which would be
				// nice and tidy, so just panic; we def don't want to continue
				panic(fmt.Sprintf("[UserUpdate] Unhandled data type (%s) for the %s flag", fl.Value.Type(), fl.Name))
			}

			if strings.HasPrefix(fl.Name, "AlertPeriod") {
				ap := &u.AlertSettings.AlertingPeriod
				setAlertingPeriodProperty(ap, strings.Replace(fl.Name, "AlertPeriod", "", -1), v)
			} else if strings.HasPrefix(fl.Name, "Alert") {
				as := &u.AlertSettings
				setAlertSettingsProperty(as, strings.Replace(fl.Name, "Alert", "", -1), v)
			} else if strings.HasPrefix(fl.Name, "Mobile") {
				ms := &u.MobileSettings
				setMobileSettingsProperty(ms, strings.Replace(fl.Name, "Mobile", "", -1), v)
			} else {
				setUserProperty(u, fl.Name, v)
			}
		}
	})

	if err := updater(); err != nil {
		return err
	}

	return nil
}
