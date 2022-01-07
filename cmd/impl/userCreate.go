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

// UserCreate is the testable implementation code for cmd.userCreateCmd
func UserCreate(f UserWriterFlags, u *api.User, creator func() error) error {
	if err := f.validate(); err != nil {
		return err
	}

	// Hydrate the user with values now known to be valid
	u.Name = f.Name
	u.Role = f.Role
	u.NotificationMethods = f.NotificationMethods
	u.MonitorGroups = f.MonitorGroups
	u.JobTitle = f.JobTitle
	u.AlertSettings = api.UserAlertSettings{
		EmailFormat: f.AlertEmailFormat,
		SkipDays:    f.AlertSkipDays,
		AlertingPeriod: api.UserAlertingPeriod{
			StartTime: f.AlertStartTime,
			EndTime:   f.AlertEndTime,
		},
		DownAlertMethods:    lookupIds(f.AlertMethodsDown, api.UserNotificationMethods),
		TroubleAlertMethods: lookupIds(f.AlertMethodsTrouble, api.UserNotificationMethods),
		UpAlertMethods:      lookupIds(f.AlertMethodsUp, api.UserNotificationMethods),
		AppLogsAlertMethods: lookupIds(f.AlertMethodsAppLogs, api.UserNotificationMethods),
		AnomalyAlertMethods: lookupIds(f.AlertMethodsAnomaly, api.UserNotificationMethods),
	}
	u.MobileSettings = api.UserMobileSettings{
		CountryCode:    f.MobileCountryCode,
		Number:         f.MobileNumber,
		SMSProviderID:  f.MobileSMSProviderID,
		CallProviderID: f.MobileCallProviderID,
	}
	u.StatusIQRole = f.StatusIQRole
	u.CloudspendRole = f.CloudspendRole

	if err := creator(); err != nil {
		return err
	}

	return nil
}
