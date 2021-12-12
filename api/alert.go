package api

import (
	"fmt"
)

// MutedAlert defines the parameters of a muting event
type MutedAlert struct {
	MutedResourceList []string `json:"muted_resource_list"`
	Category          string   `json:"category"`
	ResourceGroupList []string `json:"resource_group_list"`
	MuteTimeIso       string   `json:"mute_time_iso"`
	Reason            string   `json:"reason"`
}

// MuteOptions defines the options that can be sent to mute an option
type MuteOptions struct {
	Duration        int
	MonitorIDs      []string
	MonitorGroupIDs []string
	Category        string
	Extend          bool
	Notify          bool
}

// reqBody := map[string]interface{}{
// 	"extend_mute":         false,
// 	"muted_resource_list": []string{},
// 	"resource_group_list": []string{mg.Id},
// 	"reason":              "Deploying",
// 	"category":            "G",
// 	"notify":              true,
// }

func getMutedAlerts() (*[]MutedAlert, error) {
	return nil, nil
}

func MuteAlert(o MuteOptions) (*MutedAlert, error) {
	fmt.Println("Muting!")

	// data, _ := json.Marshal(map[string]interface{}{
	// 	"mute_time":     o.Duration,
	// 	"category":      u.EmailAddress,
	// 	"user_role":     u.Role,
	// 	"notify_medium": u.NotifyMedium,
	// 	"alert_settings": map[string]interface{}{
	// 		"email_format":       0,
	// 		"dont_alert_on_days": []int{},
	// 		"alerting_period": map[string]string{
	// 			"start_time": "00:00",
	// 			"end_time":   "00:00",
	// 		},
	// 		"down":    []int{1},
	// 		"trouble": []int{1},
	// 		"up":      []int{1},
	// 		"applogs": []int{1},
	// 		"anomaly": []int{1},
	// 	},
	// })
	// req := Request{
	// 	Endpoint: fmt.Sprintf("%s/mute_alerts", os.Getenv("API_BASE_URL")),
	// 	Method:   "PUT",
	// 	Headers: http.Header{
	// 		"Accept": {"application/json; version=2.0"},
	// 	},
	// 	Body: nil,
	// }
	// req.Headers.Set(httpHeader())
	// res, err := req.Fetch()
	// if err != nil {
	// 	return nil, err
	// }

	return nil, nil
}
