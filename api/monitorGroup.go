package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

// MonitorGroup contains the data returned from any request for monitor group
// information.
type MonitorGroup struct {
	Id                    string   `json:"group_id"`
	Name                  string   `json:"display_name"`
	Description           string   `json:"description"`
	GroupType             int      `json:"group_type"`
	Monitors              []string `json:"monitors"`
	SelectionType         int      `json:"selection_type"`
	DependencyResourceIds []string `json:"dependency_resource_ids"`
	SuppressAlert         bool     `json:"suppress_alert"`
	HealthThresholdCount  int      `json:"health_threshold_count"`
	Tags                  []string `json:"-"`
}

// getMonitorGroups returns all existing monitor groups
func getMonitorGroups() ([]MonitorGroup, error) {
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

	var groups []MonitorGroup
	err = json.Unmarshal(res.Data, &groups)
	if err != nil {
		return nil, fmt.Errorf("[getMonitorGroups] ERROR: Unable to  parse response data (%s)", err)
	}

	return groups, nil
}

// MonitorGroupExists determines whether a given user, identified by name,
// already exists in the account
func MonitorGroupExists(name string) (bool, error) {
	groups, err := getMonitorGroups()
	if err != nil {
		return false, err
	}

	for _, g := range groups {
		if strings.ToUpper(g.Name) == strings.ToUpper(name) {
			return true, nil
		}
	}

	return false, nil
}

// Create establishes a new monitor group if a group with the same name does
// not already exist
func (mg *MonitorGroup) Create() error {
	exists, err := MonitorGroupExists(mg.Name)
	if err != nil {
		return err
	}
	if exists {
		return &ConflictError{fmt.Sprintf("[MonitorGroup.Create] CONFLICTERROR: a monitor group with this name (%s) already exists on this account", mg.Name)}
	}

	// TODO: add optional data from flags
	data, _ := json.Marshal(map[string]interface{}{
		"display_name":           mg.Name,
		"health_threshold_count": 0,
	})
	req := Request{
		Endpoint: fmt.Sprintf("%s/monitor_groups", os.Getenv("API_BASE_URL")),
		Method:   "POST",
		Headers: http.Header{
			"Accept": {"application/json; version=2.1"},
		},
		Body: data,
	}
	req.Headers.Set(httpHeader())
	res, err := req.Fetch()
	if err != nil {
		return err
	}
	if res.Message != "success" {
		return fmt.Errorf("[MonitorGroup.Create] ERROR: There was a problem muting alerts for %s (%s)", mg.Name, res.Message)
	}

	if err = json.Unmarshal(res.Data, &mg); err != nil {
		return fmt.Errorf("[MonitorGroup.Create] ERROR: Unable to  parse response data (%s)", err)
	}

	return nil
}
