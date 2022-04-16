package user

import (
	"fmt"
	"strings"

	"github.com/spf13/pflag"
)

// GetAccessorFlags returns the default flagset that's passed to commands that
// need to read or identify a specific user.
func GetAccessorFlags() *pflag.FlagSet {
	accessorFlags := pflag.NewFlagSet("accessorFlags", pflag.ExitOnError)

	accessorFlags.StringP("id", "i", "", "A user identifier")
	accessorFlags.StringP("email", "e", "", "A user email address")

	return accessorFlags
}

// GetWriterFlags returns the default flagset that's passed to commands that
// write information to a user account.
func GetWriterFlags() *pflag.FlagSet {
	writerFlags := pflag.NewFlagSet("writerFlags", pflag.ExitOnError)

	writerFlags.StringP("name", "n", "Unnamed User", "Full name (first last) of the user, e.g. \"Fred Flintstone\"")
	writerFlags.IntP("role", "r", 0, "See https://www.site24x7.com/help/api/#user_constants")
	writerFlags.IntSliceP("notify-by", "N", []int{1}, "Medium by which the user will receive alerts")
	writerFlags.StringSliceP("monitor-groups", "g", []string{}, "List of monitor group identifiers to which the user should be assigned for receiving alerts")
	writerFlags.Int("alert-email-format", 1, "See https://www.site24x7.com/help/api/#alerting_constants")
	writerFlags.IntSlice("alert-skip-days", []int{}, "Days of the week on which the user should not be sent alerts: 0 (Sunday)-6 (Saturday) (default none")
	writerFlags.String("alert-start-time", "00:00", "The time of day when the user should start receiving alerts")
	writerFlags.String("alert-end-time", "00:00", "The time of day when the user should stop receiving alerts")
	writerFlags.IntSlice("alert-methods-down", []int{1}, "Preferred notification methods for down alerts")
	writerFlags.IntSlice("alert-methods-trouble", []int{1}, "Preferred notification methods for trouble alerts")
	writerFlags.IntSlice("alert-methods-up", []int{1}, "Preferred notification methods when service is restored")
	writerFlags.IntSlice("alert-methods-applogs", []int{1}, "Preferred notification methods for alerts related to application logs")
	writerFlags.IntSlice("alert-methods-anomaly", []int{1}, "Preferred notification methods for alerts when an anomaly is detected")
	writerFlags.Int("job-title", 0, "See https://www.site24x7.com/help/api/#job_title")
	writerFlags.String("mobile-country-code", "", "Country code for mobile phone number; required if voice and/or sms notifications are requested")
	writerFlags.String("mobile-phone-number", "", "Digits only; required if voice and/or sms notifications are requested")
	writerFlags.Int("mobile-sms-provider-id", 0, "See https://www.site24x7.com/help/api/#alerting_constants")
	writerFlags.Int("mobile-call-provider-id", 0, "See https://www.site24x7.com/help/api/#alerting_constants")
	writerFlags.Int("resource-type", 0, "See https://www.site24x7.com/help/api/#resource_type_constants")
	writerFlags.Int("statusiq-role", 0, "See https://www.site24x7.com/help/api/#user_constants")
	writerFlags.Int("cloudspend-role", 0, "See https://www.site24x7.com/help/api/#user_constants")
	// Not a user property, just something to pass on the request
	writerFlags.Bool("non-eu-alert-consent", false, "Mandatory for EU DC; by passing true, you confirm your consent to transfer alert-related data")

	return writerFlags
}

// lookup checks each value in a slice against the keys of a map and returns a
// slice containing the valid keys.
func lookup(keys []int, lookup map[int]string) []int {
	var result []int

	for _, i := range keys {
		if _, ok := lookup[i]; ok {
			result = append(result, i)
		}
	}

	return result
}

// validateAccessors validates user data that is passed from the command line
// specifically for the purpose of retrieving an existing user.
func validateAccessors(fs *pflag.FlagSet) error {
	i, _ := fs.GetString("id")
	e, _ := fs.GetString("email")

	if i != "" && e != "" {
		return fmt.Errorf("please include either an ID OR an email address, not both")
	} else if i == "" && e == "" {
		return fmt.Errorf("either an ID or an email address is required to identify a user")
	}

	return nil
}

// validateWriters validates writable values passed to the command via flags.
// This method only validates the input value itself, not its use or usability
// within the context of the overall upstream system. It also only validates
// flags that were changed; we should be able to safely assume that default
// values are valid.
func validateWriters(fs *pflag.FlagSet) {
	fs.Visit(func(f *pflag.Flag) {
		switch f.Name {
		case "role":
			v, _ := fs.GetInt(f.Name)
			if _, ok := RoleLookup[v]; !ok {
				panic("[validateWriters] invalid role; see https://www.site24x7.com/help/api/#user_constants")
			}
		case "job-title":
			if fs.Changed(f.Name) { // value can be nil and the flag default is nil
				t, _ := fs.GetInt("job-title")
				if _, ok := JobTitles[t]; !ok {
					panic("[validateWriters] invalid job title; see https://www.site24x7.com/help/api/#job_title")
				}
			}
		case "notify-by":
			v, _ := fs.GetIntSlice(f.Name)
			if v := lookup(v, NotificationMethods); v == nil {
				panic("[validateWriters] invalid notification method(s); see https://www.site24x7.com/help/api/#alerting_constants")
			}
		case "resource-type":
			v, _ := fs.GetInt("resource-type")
			if _, ok := ResourceTypes[v]; !ok {
				panic("[validateWriters] invalid resource type")
			}
		case "alert-email-format":
			v, _ := fs.GetInt(f.Name)
			if _, ok := EmailFormats[v]; !ok {
				panic("[validateWriters] invalid email format; see https://www.site24x7.com/help/api/#alerting_constants")
			}
		case "alert-skip-days":
			v, _ := fs.GetIntSlice(f.Name)
			if len(v) > 7 {
				panic("[validateWriters] there are only 7 days in a week")
			}
			for _, d := range v {
				if d < 0 || d > 6 {
					panic("[validateWriters] invalid skip days identified; please use 0 (Sunday) - 6 (Saturday)")
				}
			}
		case "alert-methods-down":
			d, _ := fs.GetIntSlice("alert-methods-down")
			if v := lookup(d, NotificationMethods); v == nil {
				panic("[validateWriters] invalid DOWN alert notification method(s); see https://www.site24x7.com/help/api/#alerting_constants")
			}
		case "alert-methods-trouble":
			t, _ := fs.GetIntSlice("alert-methods-trouble")
			if v := lookup(t, NotificationMethods); v == nil {
				panic("[validateWriters] invalid TROUBLE alert notification method(s); see https://www.site24x7.com/help/api/#alerting_constants")
			}
		case "alert-methods-up":
			u, _ := fs.GetIntSlice("alert-methods-up")
			if v := lookup(u, NotificationMethods); v == nil {
				panic("[validateWriters] invalid UP alert notification method(s); see https://www.site24x7.com/help/api/#alerting_constants")
			}
		case "alert-methods-applogs":
			al, _ := fs.GetIntSlice("alert-methods-applogs")
			if v := lookup(al, NotificationMethods); v == nil {
				panic("[validateWriters] invalid APPLOGS alert notification method(s); see https://www.site24x7.com/help/api/#alerting_constants")
			}
		case "alert-methods-anomaly":
			a, _ := fs.GetIntSlice("alert-methods-anomaly")
			if v := lookup(a, NotificationMethods); v == nil {
				panic("[validateWriters] invalid ANOMALY alert notification method(s); https://www.site24x7.com/help/api/#alerting_constants")
			}

		case "statusiq-role":
			v, _ := fs.GetInt(f.Name)
			if _, ok := StatusIQRoles[v]; !ok {
				panic("[validateWriters] invalid status IQ role; see https://www.site24x7.com/help/api/#user_constants")
			}
		case "cloudspend-role":
			v, _ := fs.GetInt(f.Name)
			if _, ok := CloudspendRoles[v]; !ok {
				panic("[validateWriters] invalid cloudspend role; see https://www.site24x7.com/help/api/#user_constants")
			}
		}
	})
}

// normalizeName maps a flag name to a property name
func normalizeName(f *pflag.Flag) string {
	switch f.Name {
	// for this one, the flag matches the Site24x7 terminology, but internally
	// I think "notification methods" makes more sense
	case "notify-by":
		return "NotificationMethods"

	// Handle nested properties cleanly
	case "alert-start-time":
		return "AlertingPeriodStartTime"
	case "alert-end-time":
		return "AlertingPeriodEndTime"
	case "alert-methods-down":
		return "AlertDownNotificationMethods"
	case "alert-methods-trouble":
		return "AlertTroubleNotificationMethods"
	case "alert-methods-up":
		return "AlertUpNotificationMethods"
	case "alert-methods-applogs":
		return "AlertAppLogsNotificationMethods"
	case "alert-methods-anomaly":
		return "AlertAnomalyNotificationMethods"

	// The next few cases have abbreviations ("IQ", "SMS", etc.) that we have
	// to case manually
	case "statusiq-role":
		return "StatusIQRole"
	case "mobile-sms-provider-id":
		return "MobileSMSProviderID"
	case "mobile-call-provider-id":
		return "MobileCallProviderID"

	// Everything else aligns pretty well with a "-" to CamelCase inflection
	default:
		t := strings.Title(f.Name)
		return strings.Replace(t, "-", "", -1)
	}
}
