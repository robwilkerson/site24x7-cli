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

// UserRoles maps roles available for Site24x7 access
var UserRoles = map[string]int{
	"NoAccess":        0,
	"SuperAdmin":      1,
	"Admin":           2,
	"Operator":        3,
	"BillingContact":  4,
	"SpokesPerson":    5,
	"HostingProvider": 6,
	"ReadOnly":        10,
}

// cloudspendRoles maps roles available for accessing CloudSpend
var UserCloudspendRoles = map[string]int{
	"CostAdmin": 11,
	"CostUser":  12,
}

// statusIQRoles maps roles available for accessing StatusIQ
var UserStatusIQRoles = map[string]int{
	"StatusIQSuperAdmin":     21,
	"StatusIQAdmin":          22,
	"StatusIQSpokesPerson":   23,
	"StatusIQBillingContact": 24,
	"StatusIQReadOnly":       25,
}

// notifyMediums are communication channels through which alerts can be sent
var UserNotifyMediums = map[string]int{
	"EMAIL":   1,
	"SMS":     2,
	"VOICE":   3,
	"IM":      4,
	"TWITTER": 5,
}

// emailFormats are the possible formats that can be used for notification emails
var UserEmailFormats = map[string]int{
	"Text": 0,
	"HTML": 1,
}

// jobTitles are the supported options for user job title
var UserJobTitles = map[string]int{
	"Engineer":  1,
	"CloudEng":  2,
	"DevOps":    3,
	"Webmaster": 4,
	"CLevel":    5,
	"IT":        6,
	"Others":    7,
}

// User defines the user data returned by Site24x7's user endpoints
type User struct {
	Id               string                 `json:"user_id"`
	Name             string                 `json:"display_name"`
	EmailAddress     string                 `json:"email_address"`
	Role             int                    `json:"user_role"`
	ImagePresent     bool                   `json:"image_present"`
	TwitterSettings  map[string]interface{} `json:"twitter_settings"`
	SelectionType    int                    `json:"selection_type"`
	IsAccountContact bool                   `json:"is_account_contact"`
	AlertSettings    map[string]interface{} `json:"alert_settings"`
	UserGroups       []string               `json:"user_groups"`
	IsInvited        bool                   `json:"is_invited"`
	ImSettings       map[string]interface{} `json:"im_settings"`
	NotifyMedium     []int                  `json:"notify_medium"`
	IsEditAllowed    bool                   `json:"is_edit_allowed"`
	MobileSettings   map[string]interface{} `json:"mobile_settings"`
	StatusIqRole     int                    `json:"statusiq_role"`
	CloudspendRole   int                    `json:"cloudspend_role"`
	JobTitle         int                    `json:"job_title"`
	Zuid             string                 `json:"zuid"`
}

// getUsers returns all users on the account
func getUsers() ([]User, error) {
	// Build the request
	endpoint := fmt.Sprintf("%s/users", os.Getenv("API_BASE_URL"))
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("[getUsers] ERROR: Unable to create request (%s)", err)
	}
	authH, authV := httpHeader()
	req.Header.Set("Accept", "application/json; version=2.0")
	req.Header.Set(authH, authV)

	// Send the request
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("[getUsers] ERROR: unable to execute request (%s)", err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("[getUsers] ERROR: Unable to read response body (%s)", err)
	}

	// Unmarshal the top level response
	var r ApiResponse
	if err := json.Unmarshal(body, &r); err != nil {
		return nil, fmt.Errorf("[getUsers] ERROR: Unable to  parse response body (%s)", err)
	}

	// Unmarshal the response data component
	var u []User
	err = json.Unmarshal(r.Data, &u)
	if err != nil {
		return nil, fmt.Errorf("[getUsers] ERROR: Unable to  parse response.Data (%s)", err)
	}

	return u, nil
}

// UserExists determines whether a given user, identified by email address,
// already exists in the Site24x7 account
func UserExists(email string) (bool, error) {
	users, err := getUsers()
	if err != nil {
		return false, err
	}

	for _, u := range users {
		if strings.ToUpper(u.EmailAddress) == strings.ToUpper(email) {
			return true, nil
		}
	}

	return false, nil
}

// Create creates a new Site24x7 user. This method assumes that the CLI handler
// has hydrated the user struct.
func (u *User) Create() error {
	// See whether this user already exists
	exists, err := UserExists(u.EmailAddress)
	if err != nil {
		return err
	}
	if exists {
		return &ConflictError{fmt.Sprintf("[User.Create] CONFLICTERROR: a user with this email address (%s) already exists on this account", u.EmailAddress)}
	}

	// Build the request

	// TODO: Send alert_settings config as flag options; these are hard coded
	endpoint := fmt.Sprintf("%s/users", os.Getenv("API_BASE_URL"))
	reqBody := map[string]interface{}{
		"display_name":  u.Name,
		"email_address": u.EmailAddress,
		"user_role":     u.Role,
		"notify_medium": u.NotifyMedium,
		"alert_settings": map[string]interface{}{
			"email_format":       0,
			"dont_alert_on_days": []int{},
			"alerting_period": map[string]string{
				"start_time": "00:00",
				"end_time":   "00:00",
			},
			"down":    []int{1},
			"trouble": []int{1},
			"up":      []int{1},
			"applogs": []int{1},
			"anomaly": []int{1},
		},
	}
	b, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("[User.Create] ERROR: Unable to create request body (%s)", err)
	}
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(b))
	if err != nil {
		return fmt.Errorf("[User.Create] ERROR: Unable to create request (%s)", err)
	}
	authH, authV := httpHeader()
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	req.Header.Set("Accept", "application/json; version=2.0")
	req.Header.Set(authH, authV)

	// Send the request

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("[User.Create] ERROR: unable to execute request (%s)", err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("[User.Create] ERROR: Unable to read response body (%s)", err)
	}

	// Unmarshal the top level response
	var r ApiResponse
	if err := json.Unmarshal(body, &r); err != nil {
		return fmt.Errorf("[User.Create] ERROR: Unable to parse response body (%s)", err)
	}
	// If there's no data and/or the message doesn't specify success, something went sideways
	if r.Data == nil || r.Message != "success" {
		return fmt.Errorf("[User.Create] ERROR: %s", r.Message)
	}

	// Unmarshal the domain data from the response
	err = json.Unmarshal(r.Data, &u)
	if err != nil {
		return fmt.Errorf("[User.Create] ERROR: Unable to  parse response data (%s)", err)
	}

	if r.Message == "success" {
		return nil
	} else {
		return fmt.Errorf("[User.Create] ERROR: An unexpected error occurred; code: %d, message: %s", r.Code, r.Message)
	}
}
