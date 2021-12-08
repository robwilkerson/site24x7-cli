// The api/ package provides an interface into the Site24x7 API

package api

import (
	"encoding/json"
)

// ApiResponse defines the top level schema of (almost?) every Site24x7 API
// response. The Data component contains the domain model data and varies wildly
// between API calls, so it must be flexible. We'll handle it as raw JSON at
// this stage and unmarshal it separately when we need it.
type ApiResponse struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
}

// EmptyApiResponse defines the schema of a response that includes no data. One
// example of such a response comes from the POST /milestone endpoint.
// https://www.site24x7.com/help/api/#add-a-milestone-marker
type EmptyApiResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// NotFoundError defines a custom error that should be returned when an entity
// being fetched cannot be found.
type NotFoundError struct {
	Message string
}

// Error returns a custom NotFoundError
func (e *NotFoundError) Error() string {
	return e.Message
}

// ExistsError defines a custom error that should be returned when an entity
// being created already exists.
type ConflictError struct {
	Message string
}

// Error returns a custom ConflictError
func (e *ConflictError) Error() string {
	return e.Message
}
