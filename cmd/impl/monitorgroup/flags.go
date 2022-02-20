package monitorgroup

import (
	"strings"

	"github.com/spf13/pflag"
)

// GetWriterFlags returns the default flagset that's passed to commands that
// write information to a user account.
func GetWriterFlags() *pflag.FlagSet {
	writerFlags := pflag.NewFlagSet("writerFlags", pflag.ExitOnError)

	writerFlags.StringP("description", "d", "", "Description of the monitor group")
	writerFlags.StringSliceP("monitors", "m", []string{}, "Identifiers of the monitors to be associated with the group")
	writerFlags.Int("health-threshold", 1, "Number of monitors' health that decide the group status.")
	writerFlags.StringSlice("dependent-monitors", []string{}, "Identifiers of dependent monitors")
	writerFlags.Bool("suppress-alert", false, "Suppress alert when a dependent monitor is down")

	return writerFlags
}

// normalizeName maps a flag name to a property name
func normalizeName(f *pflag.Flag) string {
	switch f.Name {
	case "health-threshold":
		return "HealthThresholdCount"

	// Everything else aligns pretty well with a "-" to CamelCase inflection
	default:
		t := strings.Title(f.Name)
		return strings.Replace(t, "-", "", -1)
	}
}
