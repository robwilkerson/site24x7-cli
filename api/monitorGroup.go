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

type MonitorGroupRequestBody struct {
	Name                 string   `json:"display_name"`
	Description          string   `json:"description"`
	Monitors             []string `json:"monitors"`
	HealthThresholdCount int      `json:"health_threshold_count"`
	DependentMonitors    []string `json:"dependency_resource_ids"`
	SuppressAlert        bool     `json:"suppress_alert"`
}

// MonitorGroup.toRequestBody performs a struct conversion
func (mg *MonitorGroup) toRequestBody() []byte {
	var b MonitorGroupRequestBody
	tmp, _ := json.Marshal(mg)
	json.Unmarshal(tmp, &b)
	body, _ := json.Marshal(b)

	return body
}

// getMonitorGroups returns all existing monitor groups
// func getMonitorGroups() ([]MonitorGroup, error) {
// 	req := Request{
// 		Endpoint: fmt.Sprintf("%s/monitor_groups", os.Getenv("API_BASE_URL")),
// 		Method:   "GET",
// 		Headers: http.Header{
// 			"Accept": {"application/json; version=2.1"},
// 		},
// 		Body: nil,
// 	}
// 	req.Headers.Set(httpHeader())
// 	res, err := req.Fetch()
// 	if err != nil {
// 		return nil, err
// 	}

// 	var groups []MonitorGroup
// 	err = json.Unmarshal(res.Data, &groups)
// 	if err != nil {
// 		return nil, fmt.Errorf("[getMonitorGroups] ERROR: Unable to  parse response data (%s)", err)
// 	}

// 	return groups, nil
// }

// MonitorGroupExists determines whether a given user, identified by name,
// already exists in the account
// func MonitorGroupExists(name string) (bool, error) {
// 	groups, err := getMonitorGroups()
// 	if err != nil {
// 		return false, err
// 	}

// 	for _, g := range groups {
// 		if strings.EqualFold(g.Name, name) {
// 			return true, nil
// 		}
// 	}

// 	return false, nil
// }

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
