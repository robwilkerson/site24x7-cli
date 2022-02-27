package usergroup

import (
	"strings"

	"github.com/spf13/pflag"
)

// GetWriterFlags returns the default flagset that's passed to commands that
// write information.
func GetWriterFlags() *pflag.FlagSet {
	writerFlags := pflag.NewFlagSet("writerFlags", pflag.ExitOnError)
	writerFlags.Int("product", 0, "Product for which the user group is being created; see https://www.site24x7.com/help/api/#product_constants")
	writerFlags.String("attribute-group-id", "", "Any attribute alert group that should be associated")
	writerFlags.StringSliceP("users", "u", []string{}, "Identifiers of any users that should be added to the group")

	return writerFlags
}

// normalizeName maps a flag name to a property name
func normalizeName(f *pflag.Flag) string {
	switch f.Name {
	case "attribute-group-id":
		return "AttributeGroup"

	// Everything else aligns pretty well with a "-" to CamelCase inflection
	default:
		t := strings.Title(f.Name)
		return strings.Replace(t, "-", "", -1)
	}
}
