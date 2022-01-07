//
// Implementation and supporting functions for the `user update` subcommand.
// Implementation functions make the code more testable
//
package impl

import (
	"errors"
	"fmt"
	"reflect"
	"site24x7/api"
)

// setName updates the user's display name, if appropriate
func setUserName(u *api.User, newValue string) {
	if newValue != "" && u.Name != newValue {
		u.Name = newValue
	}
}

func setUserProperty(u *api.User, property string, value interface{}) error {
	fmt.Printf("Setting %s to %v\n", property, value)
	ru := reflect.ValueOf(u)

	// verify that u is a pointer to a struct
	if ru.Kind() != reflect.Ptr || ru.Elem().Kind() != reflect.Struct {
		return errors.New("[setUserProperty] expected a pointer to struct")
	}
	// dereference the pointer
	ru = ru.Elem()

	// lookup the field by name
	f := ru.FieldByName(property)
	f.Set(reflect.ValueOf(value))

	return nil
}

// UserUpdate is the testable implementation code for cmd.userUpdateCmd
func UserUpdate(a UserAccessorFlags, f UserWriterFlags, u *api.User, updater func() error) error {
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

	// TODO: Merge existing data with data passed on flags
	// fmt.Println("AFTER GETTING")
	// fmt.Printf("%+v\n", u)

	// TODO: update the user object with any flag info that's different
	// setUserName(u, f.Name)
	// fields := reflect.VisibleFields(reflect.TypeOf(*u))
	fields := reflect.VisibleFields(reflect.TypeOf(f))
	// flags := reflect.ValueOf(f)
	fmt.Printf("%+v\n", flags)

	for _, field := range fields {
		// If the flag name matches the field name and the flag isn't a
		// zero value, update the user property
		flag := flags.FieldByName(field.Name)
		if flag.IsValid() && !flag.IsZero() {
			if err := setUserProperty(u, field.Name, flag.Interface()); err != nil {
				return err
			}
		}
	}

	fmt.Println("AFTER SETTING NAME")
	fmt.Printf("%+v\n", u)
	return nil

	// Hydrate the user with values now known to be valid
	// u.Name = f.Name
	// u.Role = f.Role
	// u.NotificationMethod = f.NotifyMethod
	// u.MonitorGroups = f.MonitorGroups
	// u.JobTitle = f.JobTitle
	// u.AlertSettings = map[string]interface{}{
	// 	"email_format":       f.AlertEmailFormat,
	// 	"dont_alert_on_days": f.AlertSkipDays,
	// 	"alerting_period": map[string]string{
	// 		"start_time": f.AlertStartTime,
	// 		"end_time":   f.AlertEndTime,
	// 	},
	// 	"down":    lookupIds(f.AlertMethodsDown, api.UserNotificationMethods),
	// 	"trouble": lookupIds(f.AlertMethodsTrouble, api.UserNotificationMethods),
	// 	"up":      lookupIds(f.AlertMethodsUp, api.UserNotificationMethods),
	// 	"applogs": lookupIds(f.AlertMethodsAppLogs, api.UserNotificationMethods),
	// 	"anomaly": lookupIds(f.AlertMethodsAnomaly, api.UserNotificationMethods),
	// }
	// u.MobileSettings = map[string]interface{}{
	// 	"country_code":     f.MobileCountryCode,
	// 	"mobile_number":    f.MobileNumber,
	// 	"sms_provider_id":  f.MobileSMSProviderID,
	// 	"call_provider_id": f.MobileCallProviderID,
	// }
	// u.StatusIQRole = f.StatusIQRole
	// u.CloudspendRole = f.CloudSpendRole

	// if err := updater(); err != nil {
	// 	return err
	// }

	return nil
}
