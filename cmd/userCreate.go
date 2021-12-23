//
// Implementation and supporting functions for the `user create` subcommand.
//
package cmd

import (
	"fmt"
	"site24x7/api"
)

// userCreateFlags contains the value of any flag sent to the command
type userCreateFlags struct {
	name                string // not validated locally
	role                int
	notifyMethod        []int
	statusIQRole        int
	cloudSpendRole      int
	alertEmailFormat    int
	alertSkipDays       []int
	alertStartTime      string // not validated locally
	alertEndTime        string // not validated locally
	alertMethodsDown    []int
	alertMethodsTrouble []int
	alertMethodsUp      []int
	alertMethodsAppLogs []int
	alertMethodsAnomaly []int
	jobTitle            int
}

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

// validate validates data passed to the command via flags
func (f userCreateFlags) validate() error {
	if _, ok := api.UserRoleLookup[f.role]; !ok {
		return fmt.Errorf("ERROR: Invalid role (%d)", f.role)
	}
	if v := lookupIds(f.notifyMethod, api.UserNotificationMethods); v == nil {
		return fmt.Errorf("ERROR: Invalid notification method(s) (%v)", f.notifyMethod)
	}
	// If a value was explicitly passed, error if it doesn't exist
	// 0 is the default value, a nil value, and not a valid lookup key, so we
	// should just ignore it if a zero value comes in
	if f.statusIQRole != 0 {
		if _, ok := api.UserStatusIQRoles[f.statusIQRole]; !ok {
			return fmt.Errorf("ERROR: Invalid status IQ role (%d)", f.statusIQRole)
		}
	}
	// If a value was explicitly passed, error if it doesn't exist
	// 0 is the default value, a nil value, and not a valid lookup key, so we
	// should just ignore it if a zero value comes in
	if f.cloudSpendRole != 0 { // 0 is the nil value, but also not
		if _, ok := api.UserCloudspendRoles[f.cloudSpendRole]; !ok {
			return fmt.Errorf("ERROR: Invalid cloudspend role (%d)", f.cloudSpendRole)
		}
	}
	if _, ok := api.UserEmailFormats[f.alertEmailFormat]; !ok {
		return fmt.Errorf("ERROR: Invalid email format (%d)", f.alertEmailFormat)
	}
	if f.alertSkipDays != nil {
		if len(f.alertSkipDays) > 7 {
			return fmt.Errorf("ERROR: There are 7 days in a week; %d skip days were sent", len(f.alertSkipDays))
		}
		for _, val := range f.alertSkipDays {
			if val < 0 || val > 6 {
				return fmt.Errorf("ERROR: Invalid skip days identified; please use 0 (Sunday) - 6 (Saturday)")
			}
		}
	}
	if v := lookupIds(f.alertMethodsDown, api.UserNotificationMethods); v == nil {
		return fmt.Errorf("ERROR: Invalid DOWN alert notification method(s) (%v)", f.alertMethodsDown)
	}
	if v := lookupIds(f.alertMethodsTrouble, api.UserNotificationMethods); v == nil {
		return fmt.Errorf("ERROR: Invalid TROUBLE alert notification method(s) (%v)", f.alertMethodsTrouble)
	}
	if v := lookupIds(f.alertMethodsUp, api.UserNotificationMethods); v == nil {
		return fmt.Errorf("ERROR: Invalid UP alert notification method(s) (%v)", f.alertMethodsUp)
	}
	if v := lookupIds(f.alertMethodsAppLogs, api.UserNotificationMethods); v == nil {
		return fmt.Errorf("ERROR: Invalid APPLOGS alert notification method(s) (%v)", f.alertMethodsAppLogs)
	}
	if v := lookupIds(f.alertMethodsAnomaly, api.UserNotificationMethods); v == nil {
		return fmt.Errorf("ERROR: Invalid ANOMALY alert notification method(s) (%v)", f.alertMethodsAnomaly)
	}
	if f.jobTitle != 0 { // value can be nil and the flag default is nil
		if _, ok := api.UserJobTitles[f.jobTitle]; !ok {
			return fmt.Errorf("ERROR: Invalid job title (%d)", f.jobTitle)
		}
	}

	return nil
}

// userGet is the testable implementation code for userCreateCmd
func userCreate(f userCreateFlags, u *api.User, creator func() error) error {
	if err := f.validate(); err != nil {
		return err
	}

	// Hydrate the user with values now known to be valid
	u.Name = f.name
	u.Role = f.role
	u.NotificationMethod = f.notifyMethod
	u.JobTitle = f.jobTitle
	u.AlertSettings = map[string]interface{}{
		"email_format":       f.alertEmailFormat,
		"dont_alert_on_days": f.alertSkipDays,
		"alerting_period": map[string]string{
			"start_time": f.alertStartTime,
			"end_time":   f.alertEndTime,
		},
		"down":    lookupIds(f.alertMethodsDown, api.UserNotificationMethods),
		"trouble": lookupIds(f.alertMethodsTrouble, api.UserNotificationMethods),
		"up":      lookupIds(f.alertMethodsUp, api.UserNotificationMethods),
		"applogs": lookupIds(f.alertMethodsAppLogs, api.UserNotificationMethods),
		"anomaly": lookupIds(f.alertMethodsAnomaly, api.UserNotificationMethods),
	}
	u.StatusIQRole = f.statusIQRole
	u.CloudspendRole = f.cloudSpendRole

	if err := creator(); err != nil {
		return err
	}

	return nil
}
