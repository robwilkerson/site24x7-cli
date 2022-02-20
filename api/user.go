package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"site24x7/logger"
	"strings"
)

// UserAlertingPeriod sets the window of time during which alerts may be sent
// to a user. Defaults to a 24 hour window: 00:00-00:00
type AlertingPeriod struct {
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
}

// UserAlertSettings defines a set of alert preferences
type AlertSettings struct {
	EmailFormat                int            `json:"email_format"`
	SkipDays                   []int          `json:"dont_alert_on_days"`
	AlertingPeriod             AlertingPeriod `json:"alerting_period"`
	DownNotificationMethods    []int          `json:"down"`
	TroubleNotificationMethods []int          `json:"trouble"`
	UpNotificationMethods      []int          `json:"up"`
	AppLogsNotificationMethods []int          `json:"applogs"`
	AnomalyNotificationMethods []int          `json:"anomaly"`
}

// UserMobileSettings provides details for sending alerts to a mobile device
type MobileSettings struct {
	CountryCode    string `json:"country_code"`
	PhoneNumber    string `json:"mobile_number"`
	SMSProviderID  int    `json:"sms_provider_id"`
	CallProviderID int    `json:"call_provider_id"`
}

// UserRequestBody
type UserRequestBody struct {
	Name                  string                 `json:"display_name"`
	EmailAddress          string                 `json:"email_address"`
	Role                  int                    `json:"user_role"`
	JobTitle              int                    `json:"job_title"`
	AlertSettings         AlertSettings          `json:"alert_settings"`
	MonitorGroups         []string               `json:"user_groups"`
	NotificationMethods   []int                  `json:"notify_medium"`
	MobileSettings        MobileSettings         `json:"mobile_settings"`
	StatusIQRole          int                    `json:"statusiq_role,omitempty"`
	CloudspendRole        int                    `json:"cloudspend_role,omitempty"`
	ResourceType          int                    `json:"selection_type"`
	TwitterSettings       map[string]interface{} `json:"twitter_settings,omitempty"`
	ConsentForNonEUALerts bool                   `json:"consent_for_non_eu_alerts"`
}

// User defines the user data returned by Site24x7's user endpoints
type User struct {
	ID                  string                 `json:"user_id"`
	Name                string                 `json:"display_name"`
	EmailAddress        string                 `json:"email_address"`
	Role                int                    `json:"user_role"`
	JobTitle            int                    `json:"job_title"`
	AlertSettings       AlertSettings          `json:"alert_settings"`
	MonitorGroups       []string               `json:"user_groups"`
	NotificationMethods []int                  `json:"notify_medium"`
	MobileSettings      MobileSettings         `json:"mobile_settings"`
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

// toRequestBody performs a struct conversion
func (u *User) toRequestBody() []byte {
	var b UserRequestBody
	tmp, _ := json.Marshal(u)
	json.Unmarshal(tmp, &b)
	body, _ := json.Marshal(b)

	return body
}

// UserList returns all users on the account
func UserList() (json.RawMessage, error) {
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

	if res.Message != "success" || res.Data == nil {
		return nil, fmt.Errorf("Error retrieving users; message: %s", res.Message)
	}

	return res.Data, nil
}

// UserCreate creates a new user account
func UserCreate(u *User) (json.RawMessage, error) {
	b := u.toRequestBody()

	logger.Debug(fmt.Sprintf("Request body\n%s", string(b)))

	req := Request{
		Endpoint: fmt.Sprintf("%s/users", os.Getenv("API_BASE_URL")),
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
	if res.Data == nil || res.Message != "success" {
		logger.Debug(fmt.Sprintf("Response\n%+v", res))
		if strings.HasPrefix(res.Message, "Email is already registered") {
			// Handle a "known" error just a little bit more cleanly
			return nil, &ConflictError{"a user with that email address already exists"}
		} else {
			return nil, fmt.Errorf("[User.Create] API Response error; %s", res.Message)
		}
	}

	return res.Data, nil
}

// UserGet fetches an account user
func UserGet(userID string) (json.RawMessage, error) {
	req := Request{
		Endpoint: fmt.Sprintf("%s/users/%s", os.Getenv("API_BASE_URL"), userID),
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
		return nil, &NotFoundError{"user not found"}
	}

	return res.Data, nil
}

// UserUpdate modifies an account user. https://www.site24x7.com/help/api/#update-user
func UserUpdate(u *User) (json.RawMessage, error) {
	b := u.toRequestBody()

	logger.Debug(fmt.Sprintf("Request body\n%s", string(b)))

	req := Request{
		Endpoint: fmt.Sprintf("%s/users/%s", os.Getenv("API_BASE_URL"), u.ID),
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
	if res.Data == nil || res.Message != "success" {
		return nil, fmt.Errorf("[User.Update] API Response error; %s", res.Message)
	}

	return res.Data, nil
}

// Delete removes a user from the account
func UserDelete(id string) error {
	req := Request{
		Endpoint: fmt.Sprintf("%s/users/%s", os.Getenv("API_BASE_URL"), id),
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
		return fmt.Errorf("[UserDelete] API Response error; %s", res.Message)
	}

	return nil
}
