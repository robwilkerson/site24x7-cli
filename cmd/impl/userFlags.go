//
// User types that are - or can be - shared across multiple user<Action> files
//
package impl

import (
	"fmt"
	"site24x7/api"
	"strings"

	"github.com/spf13/pflag"
)

// UserAccessorFlags define optional data points that may be sent via command
// line flags for operations that are or requre a getter operation
// (e.g. `user get`, `user delete`)
type UserAccessorFlags struct {
	ID           string
	EmailAddress string
}

// UserWriterFlags define optional data points that may be sent via command
// line flags for write operations (e.g. `user create`, `user update`)
type UserWriterFlags struct {
	Name                 string // not validated locally
	Role                 int
	NotificationMethods  []int
	MonitorGroups        []string
	AlertEmailFormat     int
	AlertSkipDays        []int
	AlertPeriodStartTime string // not validated locally
	AlertPeriodEndTime   string // not validated locally
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
	NonEUAlertConsent    bool
	ResourceType         int
	StatusIQRole         int
	CloudspendRole       int
}

// validate validates user data passed to the `user delete` command
func (f UserAccessorFlags) validate() error {
	if f.ID != "" && f.EmailAddress != "" {
		return fmt.Errorf("please include either an ID OR an email address, not both")
	} else if f.ID == "" && f.EmailAddress == "" {
		return fmt.Errorf("either an ID or an email address is required to identify a user")
	}

	return nil
}

// validate validates data passed to the command via flags. This method only
// validates the input value itself, not its use or usability within the context
// of the overall upstream system.
func (f UserWriterFlags) validate() error {
	if _, ok := api.UserRoleLookup[f.Role]; !ok {
		return fmt.Errorf("ERROR: Invalid role (%d)", f.Role)
	}
	if v := lookupIds(f.NotificationMethods, api.UserNotificationMethods); v == nil {
		return fmt.Errorf("ERROR: Invalid notification method(s) (%v)", f.NotificationMethods)
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
	if f.CloudspendRole != 0 { // 0 is the nil value, but also not
		if _, ok := api.UserCloudspendRoles[f.CloudspendRole]; !ok {
			return fmt.Errorf("ERROR: Invalid cloudspend role (%d)", f.CloudspendRole)
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
	if _, ok := api.UserResourceTypes[f.ResourceType]; !ok {
		return fmt.Errorf("ERROR: Invalid resource type (%d)", f.ResourceType)
	}

	// NOTE: There some business logic aspects that we _could_ validate, but
	// that feels like a pretty dark road. For example, in order to select a
	// text or voice call notification method, mobile settings must be sent.
	// Going to punt on that kind of thing and just focus on validating actual
	// input values.

	return nil
}

// PropertyMapper normalizes flag names to align with user properties
func PropertyMapper(f *pflag.FlagSet, name string) pflag.NormalizedName {
	switch name {
	// for this one, the flag matches the Site24x7 terminology, but internally
	// I think "notification methods" makes more sense
	case "notify-by":
		name = "NotificationMethods"

	// Handle nested properties cleanly
	case "alert-methods-down":
		name = "AlertDownNotificationMethods"
	case "alert-methods-trouble":
		name = "AlertTroubleNotificationMethods"
	case "alert-methods-up":
		name = "AlertUpNotificationMethods"
	case "alert-methods-applogs":
		name = "AlertAppLogsNotificationMethods"
	case "alert-methods-anomaly":
		name = "AlertAnomalyNotificationMethods"

	// The next few cases have abbreviations ("IQ", "SMS", etc.) that we have
	// to case manually
	case "statusiq-role":
		name = "StatusIQRole"
	case "mobile-sms-provider-id":
		name = "MobileSMSProviderID"
	case "mobile-call-provider-id":
		name = "MobileCallProviderID"

	// Everything else aligns pretty well with a "-" to CamelCase inflection
	default:
		t := strings.Title(name)
		name = strings.Replace(t, "-", "", -1)
	}

	return pflag.NormalizedName(name)
}
