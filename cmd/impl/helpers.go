package impl

import (
	"fmt"
	"reflect"
	"site24x7/logger"

	"github.com/spf13/pflag"
)

// SetProperty sets the value of a valid, addressable struct property
func SetProperty(v any, property string, value any) bool {
	t := reflect.TypeOf(v)

	// Identify the struct name for messaging clarity
	var n string
	if t.Kind() == reflect.Ptr {
		n = "*" + t.Elem().Name()
	} else {
		n = t.Name()
	}

	logger.Info(fmt.Sprintf("Setting %s.%s; value: %v\n", n, property, value))

	rv := reflect.ValueOf(v)

	if rv.Kind() != reflect.Ptr || rv.Elem().Kind() != reflect.Struct {
		logger.Warn("[impl.SetProperty] Unable to set property; v is not a pointer to a struct; ignoring")
		return false
	}

	// dereference the pointer
	rv = rv.Elem()

	// lookup the field by name and set the new value
	f := rv.FieldByName(property)

	if !f.IsValid() {
		logger.Warn(fmt.Sprintf("[impl.SetProperty] Invalid property %s; ignoring", property))
		return false
	}
	if !f.CanSet() {
		logger.Warn(fmt.Sprintf("[impl.SetProperty] Unable to set a value on an unaddressable or private field (%s); ignoring", property))
		return false
	}

	f.Set(reflect.ValueOf(value))
	return true
}

// Extracts the typed value from a flag
func TypedFlagValue(fs *pflag.FlagSet, f *pflag.Flag) any {
	var v any

	switch f.Value.Type() {
	case "string":
		v, _ = fs.GetString(f.Name)
	case "int":
		v, _ = fs.GetInt(f.Name)
	case "stringSlice":
		v, _ = fs.GetStringSlice(f.Name)
	case "intSlice":
		v, _ = fs.GetIntSlice(f.Name)
	case "bool":
		v, _ = fs.GetBool(f.Name)
	default:
		// This is a problem, but I'm not sure it needs to be a fatal one
		logger.Warn(fmt.Sprintf("[impl.TypedFlagValue] Unhandled data type (%s) for the %s flag", f.Value.Type(), f.Name))
	}

	return v
}
