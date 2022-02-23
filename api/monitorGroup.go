package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"site24x7/logger"
	"strconv"
)

// MonitorGroup contains the data returned from any request for monitor group
// information.
type MonitorGroup struct {
	ID                   string         `json:"group_id"`
	Name                 string         `json:"display_name"`
	Description          string         `json:"description"`
	Monitors             []string       `json:"monitors"`
	HealthThresholdCount int            `json:"health_threshold_count"`
	DependentMonitors    []string       `json:"dependency_resource_ids"`
	SuppressAlert        bool           `json:"suppress_alert"`
	GroupType            int            `json:"group_type,omitempty"` // https://www.site24x7.com/help/api/#monitor_group_type_constants
	Type                 int            `json:"type,omitempty"`       // https://www.site24x7.com/help/api/#monitor_group_resource_type_constants
	Tags                 []string       `json:"tags,omitempty"`
	Subgroups            []MonitorGroup `json:"subgroups,omitempty"`
}

// MonitorGroupRequestBody identifies the HTTP request body structure
type MonitorGroupRequestBody struct {
	Name                 string   `json:"display_name"`
	Description          string   `json:"description"`
	Monitors             []string `json:"monitors"`
	HealthThresholdCount int      `json:"health_threshold_count"`
	DependentMonitors    []string `json:"dependency_resource_ids"`
	SuppressAlert        bool     `json:"suppress_alert"`
}

// toRequestBody performs a struct conversion
func (mg *MonitorGroup) toRequestBody() []byte {
	var b MonitorGroupRequestBody
	tmp, _ := json.Marshal(mg)
	json.Unmarshal(tmp, &b)
	body, _ := json.Marshal(b)

	return body
}

// MonitorGroupList returns all monitor groups
func MonitorGroupList(withSubgroups bool) (json.RawMessage, error) {
	req := Request{
		Endpoint: fmt.Sprintf("%s/monitor_groups", os.Getenv("API_BASE_URL")),
		Method:   "GET",
		Headers: http.Header{
			"Accept": {"application/json; version=2.1"},
		},
		Body: nil,
		QueryString: url.Values{
			"subgroup_required": {strconv.FormatBool(withSubgroups)},
		},
	}
	req.Headers.Set(httpHeader())
	res, err := req.Fetch()
	if err != nil {
		return nil, err
	}

	if res.Message != "success" || res.Data == nil {
		return nil, fmt.Errorf("Error retrieving monitor groups; message: %s", res.Message)
	}

	return res.Data, nil
}

// MonitorGroupCreate establishes a new monitor group if a group with the same name does
// not already exist
func MonitorGroupCreate(mg *MonitorGroup) (json.RawMessage, error) {
	b := mg.toRequestBody()

	logger.Debug(fmt.Sprintf("Request body\n%s", string(b)))

	req := Request{
		Endpoint: fmt.Sprintf("%s/monitor_groups", os.Getenv("API_BASE_URL")),
		Method:   "POST",
		Headers: http.Header{
			"Accept": {"application/json; version=2.1"},
		},
		Body: b,
	}
	req.Headers.Set(httpHeader())
	res, err := req.Fetch()
	if err != nil {
		return nil, err
	}
	if res.Data == nil || res.Message != "success" {
		logger.Debug(fmt.Sprintf("Response\n%+v", res))

		return nil, fmt.Errorf("[MonitorGroup.Create] API Response error; %s", res.Message)
	}

	return res.Data, nil
}

// MonitorGroupGet fetches a monitor group
func MonitorGroupGet(id string) (json.RawMessage, error) {
	req := Request{
		Endpoint: fmt.Sprintf("%s/monitor_groups/%s", os.Getenv("API_BASE_URL"), id),
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

	if res.Data == nil {
		// Handle a "known" error just a little bit more cleanly
		return nil, &NotFoundError{"monitor group not found"}
	}

	return res.Data, nil
}

// MonitorGroupDelete removes a monitor group
func MonitorGroupDelete(id string) error {
	req := Request{
		Endpoint: fmt.Sprintf("%s/monitor_groups/%s", os.Getenv("API_BASE_URL"), id),
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
		return fmt.Errorf("[MonitorGroupDelete] API Response error; %s", res.Message)
	}

	return nil
}
