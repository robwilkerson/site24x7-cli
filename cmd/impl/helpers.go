package impl

import (
	"fmt"
	"reflect"
	"site24x7/logger"
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
