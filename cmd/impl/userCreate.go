//
// Implementation and supporting functions for the `user create` subcommand.
// Implementation functions make the code more testable
//
package impl

import (
	"site24x7/api"
)

// lookupIds checks a list of key values against a map and returns a slice
// containing each existing key's value
func lookupIds(keys []int, lookup map[int]string) []int {
	var result []int

	for _, i := range keys {
		if _, ok := lookup[i]; ok {
			result = append(result, i)
		}
	}

	return result
}

// UserGet is the testable implementation code for cmd.userCreateCmd
func UserCreate(f UserWriterFlags, u *api.User, creator func() error) error {
	if err := f.validate(); err != nil {
		return err
	}

	// Hydrate the user with values now known to be valid
	u.Name = f.Name
	u.Role = f.Role
	u.NotificationMethod = f.NotifyMethod
	u.MonitorGroups = f.MonitorGroups
	u.JobTitle = f.JobTitle
	u.AlertSettings = map[string]interface{}{
		"email_format":       f.AlertEmailFormat,
		"dont_alert_on_days": f.AlertSkipDays,
		"alerting_period": map[string]string{
			"start_time": f.AlertStartTime,
			"end_time":   f.AlertEndTime,
		},
		"down":    lookupIds(f.AlertMethodsDown, api.UserNotificationMethods),
		"trouble": lookupIds(f.AlertMethodsTrouble, api.UserNotificationMethods),
		"up":      lookupIds(f.AlertMethodsUp, api.UserNotificationMethods),
		"applogs": lookupIds(f.AlertMethodsAppLogs, api.UserNotificationMethods),
		"anomaly": lookupIds(f.AlertMethodsAnomaly, api.UserNotificationMethods),
	}
	u.MobileSettings = map[string]interface{}{
		"country_code":     f.MobileCountryCode,
		"mobile_number":    f.MobileNumber,
		"sms_provider_id":  f.MobileSMSProviderID,
		"call_provider_id": f.MobileCallProviderID,
	}
	u.StatusIQRole = f.StatusIQRole
	u.CloudspendRole = f.CloudSpendRole

	if err := creator(); err != nil {
		return err
	}

	return nil
}
