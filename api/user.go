package api

import (
	"encoding/json"
	"fmt"
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

// getUsers returns all users on the account.
func getUsers() ([]User, error) {
	req := Request{
		Endpoint: fmt.Sprintf("%s/users", os.Getenv("API_BASE_URL")),
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

	var u []User
	err = json.Unmarshal(res.Data, &u)
	if err != nil {
		return nil, fmt.Errorf("[getUsers] ERROR: Unable to  parse response data (%s)", err)
	}

	return u, nil
}

// UserExists determines whether a given user, identified by email address,
// already exists in the Site24x7 account.
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

// findByEmail returns a user found with a matching email address
func (u *User) findUserByEmail() error {
	users, err := getUsers()
	if err != nil {
		return err
	}

	for _, usr := range users {
		if strings.ToLower(usr.EmailAddress) == strings.ToLower(u.EmailAddress) {
			// Update the receiver with the official, fully hydrated user
			*u = usr

			return nil
		}
	}

	return &NotFoundError{fmt.Sprintf("[User.findByEmail] NOTFOUNDERROR: No user with that email (%s) found", u.EmailAddress)}
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

	// TODO: replace hard-coded values with flag data
	data, _ := json.Marshal(map[string]interface{}{
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
	})
	req := Request{
		Endpoint: fmt.Sprintf("%s/users", os.Getenv("API_BASE_URL")),
		Method:   "POST",
		Headers: http.Header{
			"Accept": {"application/json; version=2.0"},
		},
		Body: data,
	}
	req.Headers.Set(httpHeader())
	res, err := req.Fetch()
	if err != nil {
		return err
	}
	if res.Data == nil || res.Message != "success" {
		return fmt.Errorf("[User.Create] ERROR: %s", res.Message)
	}

	// Unmarshal the domain data from the response
	err = json.Unmarshal(res.Data, &u)
	if err != nil {
		return fmt.Errorf("[User.Create] ERROR: Unable to  parse response data (%s)", err)
	}

	return nil
}

// Get returns a user on the account
func (u *User) Get() error {
	// If an email address is sent, convert that to an id
	if u.EmailAddress != "" {
		if err := u.findUserByEmail(); err != nil {
			return err
		}

		return nil
	}

	// If we dropped here, we can assume that an identifier was passed. We could
	// do a "find by" operation, but getting exactly what we need should be
	// faster.
	req := Request{
		Endpoint: fmt.Sprintf("%s/users/%s", os.Getenv("API_BASE_URL"), u.Id),
		Method:   "GET",
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

	if err = json.Unmarshal(res.Data, &u); err != nil {
		return fmt.Errorf("[getUsers] ERROR: Unable to  parse response data (%s)", err)
	}

	return nil
}
