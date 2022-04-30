package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"site24x7/logger"
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

// UserGroupRequestBody defines the HTTP request body structure
type UserGroupRequestBody struct {
	Name           string   `json:"display_name"`
	Product        int      `json:"product_id"` // https://www.site24x7.com/help/api/#product_constants
	Users          []string `json:"users"`
	AttributeGroup string   `json:"attribute_group_id"`
}

// toRequestBody performs a struct conversion
func (ug *UserGroup) toRequestBody() []byte {
	var b UserGroupRequestBody
	tmp, _ := json.Marshal(ug)
	json.Unmarshal(tmp, &b)
	body, _ := json.Marshal(b)

	return body
}

// UserGroupCreate establishes a new user group
// https://www.site24x7.com/help/api/#create-user-group
func UserGroupCreate(ug *UserGroup) (json.RawMessage, error) {
	b := ug.toRequestBody()

	req := Request{
		Endpoint: fmt.Sprintf("%s/user_groups", os.Getenv("API_BASE_URL")),
		Method:   "POST",
		Headers: http.Header{
			"Accept": {"application/json; version=2.0"},
		},
		Body: b,
	}
	req.Headers.Set(httpHeader())
	res, err := req.Fetch()
	if err != nil {
		return nil, err
	}
	if string(res.Data) == "{}" || res.Message != "success" {
		logger.Debug(fmt.Sprintf("Response\n%+v", res))

		return nil, fmt.Errorf("[api.UserGroupCreate] API Response error; %s", res.Message)
	}

	return res.Data, nil
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

// UserGroupUpdate updates a user group
// https://www.site24x7.com/help/api/#update-user-group
func UserGroupUpdate(ug *UserGroup) (json.RawMessage, error) {
	b := ug.toRequestBody()

	req := Request{
		Endpoint: fmt.Sprintf("%s/user_groups/%s", os.Getenv("API_BASE_URL"), ug.ID),
		Method:   "PUT",
		Headers: http.Header{
			"Accept": {"application/json; version=2.0"},
		},
		Body: b,
	}
	req.Headers.Set(httpHeader())
	res, err := req.Fetch()
	if err != nil {
		return nil, err
	}
	if string(res.Data) == "{}" || res.Message != "success" {
		return nil, fmt.Errorf("[api.UserGroupUpdate] API Response error; %s", res.Message)
	}

	return res.Data, nil
}

// UserGroupDelete removes a user group
func UserGroupDelete(id string) error {
	req := Request{
		Endpoint: fmt.Sprintf("%s/user_groups/%s", os.Getenv("API_BASE_URL"), id),
		Method:   "DELETE",
		Headers: http.Header{
			"Accept": {"application/json; version=2.0"},
		},
		Body: nil,
	}
	req.Headers.Set(httpHeader())
	res, err := req.Fetch()
	if err != nil {
		return err
	}
	if res.Message != "success" {
		return fmt.Errorf("[api.UserGroupDelete] API Response error; %s", res.Message)
	}

	return nil
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
