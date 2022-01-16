package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

// UserRoles maps role ids to friendly names
// https://www.site24x7.com/help/api/#user_constants
var UserRoleLookup = map[int]string{
	0: "No Access",
	1: "Super Administrator",
	2: "Administrator",
	3: "Operator",
	4: "Billing Contact",
	5: "Spokesperson",
	6: "Hosting Provider",
	7: "Read Only",
}

// UserCloudspendRoles maps role ids to friendly names
// https://www.site24x7.com/help/api/#user_constants
var UserCloudspendRoles = map[int]string{
	11: "Cost Administrator",
	12: "Cost User",
}

// UserStatusIQRoles maps role ids to friendly names
// https://www.site24x7.com/help/api/#alerting_constants
var UserStatusIQRoles = map[int]string{
	21: "StatusIQ Super Administrator",
	22: "StatusIQ Administrator",
	23: "StatusIQ SpokesPerson",
	24: "StatusIQ Billing Contact",
	25: "StatusIQ Read Only",
}

// UserNotificationMethods are communication channels through which alerts can be sent
// https://www.site24x7.com/help/api/#alerting_constants
var UserNotificationMethods = map[int]string{
	1: "Email",
	2: "SMS",
	3: "Voice Call",
	4: "IM",
	5: "Twitter",
}

// UserEmailFormats are the possible formats that can be used for notification emails
// https://www.site24x7.com/help/api/#alerting_constants
var UserEmailFormats = map[int]string{
	0: "TEXT",
	1: "HTML",
}

// UserJobTitles are the supported options for user job title
// https://www.site24x7.com/help/api/#job_title
var UserJobTitles = map[int]string{
	1: "IT Engineer",
	2: "Cloud Engineer",
	3: "DevOps Engineer",
	4: "Webmaster",
	5: "CEO/CTO",
	6: "Internal IT",
	7: "Others",
}

// UserResourceTypes are the supported resource/selection types
// https://www.site24x7.com/help/api/#resource_type_constants
var UserResourceTypes = map[int]string{
	0: "All Monitors",
	1: "Monitor Group",
	2: "Monitor",
	3: "Tags",
	4: "Monitor Type",
}

// UserAlertPeriod sets the window of time during which alerts may be sent
// to a user. Defaults to a 24 hour window: 00:00-00:00
type UserAlertingPeriod struct {
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
}

// UserAlertSettings defines a set of alert preferences
type UserAlertSettings struct {
	EmailFormat                int                `json:"email_format"`
	SkipDays                   []int              `json:"dont_alert_on_days"`
	AlertingPeriod             UserAlertingPeriod `json:"alerting_period"`
	DownNotificationMethods    []int              `json:"down"`
	TroubleNotificationMethods []int              `json:"trouble"`
	UpNotificationMethods      []int              `json:"up"`
	AppLogsNotificationMethods []int              `json:"applogs"`
	AnomalyNotificationMethods []int              `json:"anomaly"`
}

// UserMobileSettings provides details for sending alerts to a mobile device
type UserMobileSettings struct {
	CountryCode    string `json:"country_code"`
	Number         string `json:"mobile_number"`
	SMSProviderID  int    `json:"sms_provider_id"`
	CallProviderID int    `json:"call_provider_id"`
}

// User defines the user data returned by Site24x7's user endpoints
type User struct {
	// TODO: "Id" --> "ID"
	Id                  string                 `json:"user_id"`
	Name                string                 `json:"display_name"`
	EmailAddress        string                 `json:"email_address"`
	Role                int                    `json:"user_role"`
	JobTitle            int                    `json:"job_title"`
	AlertSettings       UserAlertSettings      `json:"alert_settings"`
	MonitorGroups       []string               `json:"user_groups"`
	NotificationMethods []int                  `json:"notify_medium"`
	MobileSettings      UserMobileSettings     `json:"mobile_settings"`
	StatusIQRole        int                    `json:"statusiq_role"`
	CloudspendRole      int                    `json:"cloudspend_role"`
	ResourceType        int                    `json:"selection_type"`
	ImagePresent        bool                   `json:"image_present"`
	TwitterSettings     map[string]interface{} `json:"twitter_settings"`
	IsAccountContact    bool                   `json:"is_account_contact"`
	IsInvited           bool                   `json:"is_invited"`
	ImSettings          map[string]interface{} `json:"im_settings"`
	IsEditAllowed       bool                   `json:"is_edit_allowed"`
	Zuid                string                 `json:"zuid"`
}

// findByEmail returns a user found with a matching email address
func (u *User) findUserByEmail() error {
	users, err := GetUsers()
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

// getUsers returns all users on the account.
func GetUsers() ([]User, error) {
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
	users, err := GetUsers()
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

// Create creates a new Site24x7 user and hydrates a pointer
func (u *User) Create() error {
	// See whether this user already exists
	exists, err := UserExists(u.EmailAddress)
	if err != nil {
		return err
	}
	if exists {
		return &ConflictError{fmt.Sprintf("[User.Create] Conflict; a user with this email address (%s) already exists on this account", u.EmailAddress)}
	}

	// TODO: include optional data from flags
	data := map[string]interface{}{
		"display_name":    u.Name,
		"email_address":   u.EmailAddress,
		"user_role":       u.Role,
		"notify_medium":   u.NotificationMethods,
		"alert_settings":  u.AlertSettings,
		"job_title":       u.JobTitle,
		"mobile_settings": u.MobileSettings,
	}

	// 0 is the default status iq and cloudspend role, but it's not a valid
	// role for either and the call will error if sent as such. Only send them
	// if the user entered a non-default value.

	if u.StatusIQRole != 0 {
		data["statusiq_role"] = u.StatusIQRole
	}
	if u.CloudspendRole != 0 {
		data["cloudspend_role"] = u.CloudspendRole
	}

	body, _ := json.Marshal(data)

	// TODO: apply a verbose context for debug/info output?
	// fmt.Println(string(data))

	req := Request{
		Endpoint: fmt.Sprintf("%s/users", os.Getenv("API_BASE_URL")),
		Method:   "POST",
		Headers: http.Header{
			"Accept": {"application/json; version=2.0"},
		},
		Body: body,
	}
	req.Headers.Set(httpHeader())
	res, err := req.Fetch()
	if err != nil {
		return err
	}
	if res.Data == nil || res.Message != "success" {
		// fmt.Printf("%+v", res)
		return fmt.Errorf("[User.Create] API Response error; %s", res.Message)
	}

	// Unmarshal the domain data from the response
	if err = json.Unmarshal(res.Data, &u); err != nil {
		return fmt.Errorf("[User.Create] Unable to  parse response data (%s)", err)
	}

	return nil
}

// Get fetches an account user and updates a pointer
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
		return fmt.Errorf("[User.Get] Unable to  parse response data (%s)", err)
	}

	return nil
}

// Update modifies an account user. https://www.site24x7.com/help/api/#update-user
func (u *User) Update() error {
	// See whether this user already exists
	//
	exists, err := UserExists(u.EmailAddress)
	if err != nil {
		return err
	}
	if !exists {
		return &NotFoundError{"[User.Update] User not found"}
	}

	data := map[string]interface{}{
		"display_name":    u.Name,
		"email_address":   u.EmailAddress,
		"user_role":       u.Role,
		"notify_medium":   u.NotificationMethods,
		"alert_settings":  u.AlertSettings,
		"job_title":       u.JobTitle,
		"mobile_settings": u.MobileSettings,
		"user_groups":     u.MonitorGroups,
	}

	// 0 is the default status iq and cloudspend role, but it's not a valid
	// role for either and the call will error if sent as such. Only send them
	// if the user entered a non-default value.

	if u.StatusIQRole != 0 {
		data["statusiq_role"] = u.StatusIQRole
	}
	if u.CloudspendRole != 0 {
		data["cloudspend_role"] = u.CloudspendRole
	}

	body, _ := json.Marshal(data)

	// TODO: apply a verbose context for debug/info output?
	// fmt.Println(string(data))

	req := Request{
		Endpoint: fmt.Sprintf("%s/users/%s", os.Getenv("API_BASE_URL"), u.Id),
		Method:   "PUT",
		Headers: http.Header{
			"Accept": {"application/json; version=2.0"},
		},
		Body: body,
	}
	req.Headers.Set(httpHeader())
	res, err := req.Fetch()
	if err != nil {
		return err
	}
	// fmt.Printf("%+v", res)
	if res.Data == nil || res.Message != "success" {
		return fmt.Errorf("[User.Update] API Response error; %s", res.Message)
	}

	// Unmarshal the domain data from the response
	if err = json.Unmarshal(res.Data, &u); err != nil {
		return fmt.Errorf("[User.Update] Unable to  parse response data (%s)", err)
	}

	return nil
}

// Delete removes a user from the account
func (u *User) Delete() error {
	// If an email address is sent, convert that to an id
	if u.EmailAddress != "" {
		if err := u.findUserByEmail(); err != nil {
			return err
		}
	}

	req := Request{
		Endpoint: fmt.Sprintf("%s/users/%s", os.Getenv("API_BASE_URL"), u.Id),
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
		return fmt.Errorf("[User.Delete] API Response error; %s", res.Message)
	}

	return nil
}
