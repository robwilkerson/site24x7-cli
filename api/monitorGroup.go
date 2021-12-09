package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
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
	// Build the request

	endpoint := fmt.Sprintf("%s/monitor_groups", os.Getenv("API_BASE_URL"))
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("[getMonitorGroups] ERROR: Unable to create request (%s)", err)
	}
	authH, authV := httpHeader()
	req.Header.Set("Accept", "application/json; version=2.1")
	req.Header.Set(authH, authV)

	// Send the request

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("[getMonitorGroups] ERROR: unable to execute request (%s)", err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("[getMonitorGroups] ERROR: Unable to read response body (%s)", err)
	}

	// Unmarshal the top level response

	var r ApiResponse
	if err := json.Unmarshal(body, &r); err != nil {
		return nil, fmt.Errorf("[getMonitorGroups] ERROR: Unable to  parse response body (%s)", err)
	}

	// Unmarshal the response data

	var groups []MonitorGroup
	err = json.Unmarshal(r.Data, &groups)
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

	// Build the request to create a new group
	endpoint := fmt.Sprintf("%s/monitor_groups", os.Getenv("API_BASE_URL"))
	reqBody := map[string]interface{}{
		"display_name":           mg.Name,
		"health_threshold_count": 0,
	}
	b, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("[MonitorGroup.Create] ERROR: Unable to create request body (%s)", err)
	}
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(b))
	if err != nil {
		return fmt.Errorf("[MonitorGroup.Create] ERROR: Unable to create request (%s)", err)
	}
	authH, authV := httpHeader()
	req.Header.Set("Content-Type", "application/json;charset=UTF-")
	req.Header.Set("Accept", "application/json; version=2.1")
	req.Header.Set(authH, authV)

	// Send the request
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("[MonitorGroup.Create] ERROR: unable to execute request (%s)", err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("[MonitorGroup.Create] ERROR: Unable to read response body (%s)", err)
	}

	// Unmarshal the top level response
	var r ApiResponse
	if err := json.Unmarshal(body, &r); err != nil {
		return fmt.Errorf("[MonitorGroup.Create] ERROR: Unable to unmarshal response body (%s)", err)
	}

	// Unmarshal the response data component
	err = json.Unmarshal(r.Data, &mg)
	if err != nil {
		return fmt.Errorf("[MonitorGroup.Create] ERROR: Unable to  parse response data (%s)", err)
	}

	if r.Message == "success" {
		return nil
	} else {
		return fmt.Errorf("[MonitorGroup.Create] ERROR: There was a problem muting alerts for %s (%s)", mg.Name, r.Message)
	}
}
