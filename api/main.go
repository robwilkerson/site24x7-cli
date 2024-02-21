// The api/ package provides an interface into the Site24x7 API

package api

// NotFoundError defines a custom error that should be returned when an entity
// being fetched cannot be found.
type NotFoundError struct {
	Message string
}

// Error returns a custom NotFoundError
func (e *NotFoundError) Error() string {
	return e.Message
}

// ConflictError defines a custom error that should be returned when an entity
// being created already exists.
type ConflictError struct {
	Message string
}

// Error returns a custom ConflictError
func (e *ConflictError) Error() string {
	return e.Message
}
