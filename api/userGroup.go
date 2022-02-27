package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

// UserGroup contains the data returned from any request for user group
// information.
type UserGroup struct {
	ID             string   `json:"user_group_id"`
	Name           string   `json:"display_name"`
	Product        int      `json:"product_id"` // https://www.site24x7.com/help/api/#product_constants
	Users          []string `json:"users"`
	AttributeGroup string   `json:"attribute_group_id"`
}

// UserGroupGet fetches a monitor group
// https://www.site24x7.com/help/api/#retrieve-user-group
func UserGroupGet(id string) (json.RawMessage, error) {
	req := Request{
		Endpoint: fmt.Sprintf("%s/user_groups/%s", os.Getenv("API_BASE_URL"), id),
		Method:   "GET",
		Headers: http.Header{
			"Accept": {"application/json; version=2.0"},
		},
		Body: nil,
	}
	req.Headers.Set(httpHeader())
	res, err := req.Fetch()
	if err != nil {
		return nil, err
	}

	if string(res.Data) == "{}" {
		// Handle a "known" error just a little bit more cleanly
		return nil, &NotFoundError{"user group not found"}
	}

	return res.Data, nil
}

// UserGroupList returns all monitor groups
// https://www.site24x7.com/help/api/#list-of-all-user-groups
func UserGroupList() (json.RawMessage, error) {
	req := Request{
		Endpoint: fmt.Sprintf("%s/user_groups", os.Getenv("API_BASE_URL")),
		Method:   "GET",
		Headers: http.Header{
			"Accept": {"application/json; version=2.0"},
		},
		Body: nil,
	}
	req.Headers.Set(httpHeader())
	res, err := req.Fetch()
	if err != nil {
		return nil, err
	}

	if res.Message != "success" || string(res.Data) == "{}" {
		return nil, fmt.Errorf("Error retrieving user groups; message: %s", res.Message)
	}

	return res.Data, nil
}
