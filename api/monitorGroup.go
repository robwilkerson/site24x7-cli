package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"site24x7/logger"
)

// MonitorGroup contains the data returned from any request for monitor group
// information.
type MonitorGroup struct {
	ID                   string   `json:"group_id"`
	Name                 string   `json:"display_name"`
	Description          string   `json:"description"`
	Monitors             []string `json:"monitors"`
	HealthThresholdCount int      `json:"health_threshold_count"`
	DependentMonitors    []string `json:"dependency_resource_ids"`
	SuppressAlert        bool     `json:"suppress_alert"`
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
func MonitorGroupList() (json.RawMessage, error) {
	req := Request{
		Endpoint: fmt.Sprintf("%s/monitor_groups", os.Getenv("API_BASE_URL")),
		Method:   "GET",
		Headers: http.Header{
			"Accept": {"application/json; version=2.1"},
		},
		Body: nil,
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

// Create establishes a new monitor group if a group with the same name does
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

// UserGet fetches an account user
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
