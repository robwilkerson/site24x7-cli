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
	"NOACCESS":     0,
	"SUPERADMIN":   1,
	"ADMIN":        2,
	"OPERATOR":     3,
	"BILLING":      4,
	"SPOKESPERSON": 5,
	"HOSTING":      6,
	"READONLY":     10,
}

// UserCloudspendRoles maps roles available for accessing CloudSpend
var UserCloudspendRoles = map[string]int{
	"ADMIN": 11,
	"USER":  12,
}

// UserStatusIQRoles maps roles available for accessing StatusIQ
var UserStatusIQRoles = map[string]int{
	"SUPERADMIN":   21,
	"ADMIN":        22,
	"SPOKESPERSON": 23,
	"BILLING":      24,
	"READONLY":     25,
}

// UserNotifyMediums are communication channels through which alerts can be sent
var UserNotifyMediums = map[string]int{
	"EMAIL":   1,
	"SMS":     2,
	"VOICE":   3,
	"IM":      4,
	"TWITTER": 5,
}

// UserEmailFormats are the possible formats that can be used for notification emails
var UserEmailFormats = map[string]int{
	"TEXT": 0,
	"HTML": 1,
}

// UserJobTitles are the supported options for user job title
var UserJobTitles = map[string]int{
	"ENGINEER":  1,
	"CLOUDENG":  2,
	"DEVOPS":    3,
	"WEBMASTER": 4,
	"CLEVEL":    5,
	"IT":        6,
	"OTHERS":    7,
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
	StatusIQRole     int                    `json:"statusiq_role"`
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
	if err = json.Unmarshal(res.Data, &u); err != nil {
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
		if strings.EqualFold(u.EmailAddress, email) {
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
		if strings.EqualFold(usr.EmailAddress, u.EmailAddress) {
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
	// TODO: include optional data from flags
	data, _ := json.Marshal(map[string]interface{}{
		"display_name":    u.Name,
		"email_address":   u.EmailAddress,
		"user_role":       u.Role,
		"statusiq_role":   u.StatusIQRole,
		"cloudspend_role": u.CloudspendRole,
		"notify_medium":   u.NotifyMedium,
		"alert_settings":  u.AlertSettings,
		"job_title":       u.JobTitle,
	})
	fmt.Println(string(data))

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
		fmt.Printf("%+v", res)
		return fmt.Errorf("[User.Create] ERROR: %s", res.Message)
	}

	// Unmarshal the domain data from the response
	if err = json.Unmarshal(res.Data, &u); err != nil {
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
