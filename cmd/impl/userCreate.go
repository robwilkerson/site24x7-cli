//
// Implementation and supporting functions for the `user create` subcommand.
// Implementation functions make the code more testable
//
package impl

import (
	"fmt"
	"site24x7/api"
)

// userCreateFlags contains the value of any flag sent to the command
type UserCreateFlags struct {
	Name                 string // not validated locally
	Role                 int
	NotifyMethod         []int
	MonitorGroups        []string
	NonEUAlertConsent    bool
	AlertEmailFormat     int
	AlertSkipDays        []int
	AlertStartTime       string // not validated locally
	AlertEndTime         string // not validated locally
	AlertMethodsDown     []int
	AlertMethodsTrouble  []int
	AlertMethodsUp       []int
	AlertMethodsAppLogs  []int
	AlertMethodsAnomaly  []int
	JobTitle             int
	MobileCountryCode    string
	MobileNumber         string
	MobileSMSProviderID  int
	MobileCallProviderID int
	StatusIQRole         int
	CloudSpendRole       int
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

// validate validates data passed to the command via flags. This method only
// validates the input value itself, not its use or usability within the context
// of the overall upstream system.
func (f UserCreateFlags) validate() error {
	if _, ok := api.UserRoleLookup[f.Role]; !ok {
		return fmt.Errorf("ERROR: Invalid role (%d)", f.Role)
	}
	if v := lookupIds(f.NotifyMethod, api.UserNotificationMethods); v == nil {
		return fmt.Errorf("ERROR: Invalid notification method(s) (%v)", f.NotifyMethod)
	}
	// If a value was explicitly passed, error if it doesn't exist
	// 0 is the default value, a nil value, and not a valid lookup key, so we
	// should just ignore it if a zero value comes in
	if f.StatusIQRole != 0 {
		if _, ok := api.UserStatusIQRoles[f.StatusIQRole]; !ok {
			return fmt.Errorf("ERROR: Invalid status IQ role (%d)", f.StatusIQRole)
		}
	}
	// If a value was explicitly passed, error if it doesn't exist
	// 0 is the default value, a nil value, and not a valid lookup key, so we
	// should just ignore it if a zero value comes in
	if f.CloudSpendRole != 0 { // 0 is the nil value, but also not
		if _, ok := api.UserCloudspendRoles[f.CloudSpendRole]; !ok {
			return fmt.Errorf("ERROR: Invalid cloudspend role (%d)", f.CloudSpendRole)
		}
	}
	if _, ok := api.UserEmailFormats[f.AlertEmailFormat]; !ok {
		return fmt.Errorf("ERROR: Invalid email format (%d)", f.AlertEmailFormat)
	}
	if f.AlertSkipDays != nil {
		if len(f.AlertSkipDays) > 7 {
			return fmt.Errorf("ERROR: There are 7 days in a week; %d skip days were sent", len(f.AlertSkipDays))
		}
		for _, val := range f.AlertSkipDays {
			if val < 0 || val > 6 {
				return fmt.Errorf("ERROR: Invalid skip days identified; please use 0 (Sunday) - 6 (Saturday)")
			}
		}
	}
	if v := lookupIds(f.AlertMethodsDown, api.UserNotificationMethods); v == nil {
		return fmt.Errorf("ERROR: Invalid DOWN alert notification method(s) (%v)", f.AlertMethodsDown)
	}
	if v := lookupIds(f.AlertMethodsTrouble, api.UserNotificationMethods); v == nil {
		return fmt.Errorf("ERROR: Invalid TROUBLE alert notification method(s) (%v)", f.AlertMethodsTrouble)
	}
	if v := lookupIds(f.AlertMethodsUp, api.UserNotificationMethods); v == nil {
		return fmt.Errorf("ERROR: Invalid UP alert notification method(s) (%v)", f.AlertMethodsUp)
	}
	if v := lookupIds(f.AlertMethodsAppLogs, api.UserNotificationMethods); v == nil {
		return fmt.Errorf("ERROR: Invalid APPLOGS alert notification method(s) (%v)", f.AlertMethodsAppLogs)
	}
	if v := lookupIds(f.AlertMethodsAnomaly, api.UserNotificationMethods); v == nil {
		return fmt.Errorf("ERROR: Invalid ANOMALY alert notification method(s) (%v)", f.AlertMethodsAnomaly)
	}
	if f.JobTitle != 0 { // value can be nil and the flag default is nil
		if _, ok := api.UserJobTitles[f.JobTitle]; !ok {
			return fmt.Errorf("ERROR: Invalid job title (%d)", f.JobTitle)
		}
	}

	// NOTE: There some business logic aspects that we _could_ validate, but
	// that feels like a pretty dark road. For example, in order to select a
	// text or voice call notification method, mobile settings must be sent.
	// going to punt on that kind of thing and just focus on validating actual
	// input values.

	return nil
}

// userGet is the testable implementation code for userCreateCmd
func UserCreate(f UserCreateFlags, u *api.User, creator func() error) error {
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
