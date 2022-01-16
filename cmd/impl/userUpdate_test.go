//
// Implementation and supporting functions for the `user update` subcommand.
// Implementation functions make the code more testable
//
package impl

import (
	"reflect"
	"site24x7/api"
	"testing"
)

var defaultAlertMethods []int = []int{1}
var defaultAlertingPeriod api.UserAlertingPeriod = api.UserAlertingPeriod{
	StartTime: "00:00",
	EndTime:   "00:00",
}
var defaultAlertSettings api.UserAlertSettings = api.UserAlertSettings{
	EmailFormat:                1,
	SkipDays:                   []int{},
	AlertingPeriod:             defaultAlertingPeriod,
	DownNotificationMethods:    defaultAlertMethods,
	TroubleNotificationMethods: defaultAlertMethods,
	UpNotificationMethods:      defaultAlertMethods,
	AppLogsNotificationMethods: defaultAlertMethods,
	AnomalyNotificationMethods: defaultAlertMethods,
}
var defaultMobileSettings api.UserMobileSettings = api.UserMobileSettings{
	CountryCode:    "",
	Number:         "",
	SMSProviderID:  0,
	CallProviderID: 0,
}
var defaultUser *api.User = &api.User{
	Name:                "Unnamed User",
	EmailAddress:        "oompa@loompa.com",
	Role:                0,
	NotificationMethods: []int{1},
	AlertSettings:       defaultAlertSettings,
	JobTitle:            0,
	MobileSettings:      defaultMobileSettings,
	ResourceType:        0,
}

func Test_setAlertingPeriodProperty(t *testing.T) {
	type args struct {
		ap       *api.UserAlertingPeriod
		property string
		value    interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Sets the StartTime property",
			args: args{
				ap:       &defaultAlertingPeriod,
				property: "StartTime",
				value:    "04:31",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setAlertingPeriodProperty(tt.args.ap, tt.args.property, tt.args.value)
			if tt.args.property == "StartTime" && tt.args.ap.StartTime != tt.args.value {
				t.Errorf("setAlertingPeriodProperty() StartTime = %v, wantErr %v", tt.args.ap.StartTime, tt.args.value)
			}
			if tt.args.property == "EndTime" && tt.args.ap.EndTime != tt.args.value {
				t.Errorf("setAlertingPeriodProperty() EndTime = %v, wantErr %v", tt.args.ap.EndTime, tt.args.value)
			}
		})
	}
}

func Test_setAlertSettingsProperty(t *testing.T) {
	type args struct {
		as       *api.UserAlertSettings
		property string
		value    interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Updates skip days",
			args: args{
				as:       &defaultAlertSettings,
				property: "SkipDays",
				value:    []int{1, 3, 5},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setAlertSettingsProperty(tt.args.as, tt.args.property, tt.args.value)
			if tt.args.property == "SkipDays" && !reflect.DeepEqual(tt.args.as.SkipDays, tt.args.value) {
				t.Errorf("setAlertingSettingsProperty() SkipDays = %v, wantErr %v", tt.args.as.SkipDays, tt.args.value)
			}
		})
	}
}

func Test_setMobileSettingsProperty(t *testing.T) {
	type args struct {
		ms       *api.UserMobileSettings
		property string
		value    interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Updates skip days",
			args: args{
				ms:       &defaultMobileSettings,
				property: "CountryCode",
				value:    "19",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setMobileSettingsProperty(tt.args.ms, tt.args.property, tt.args.value)
			if tt.args.property == "CountryCode" && tt.args.ms.CountryCode != tt.args.value {
				t.Errorf("setMobileSettingsProperty() CountryCode = %v, expected = %v", tt.args.ms.CountryCode, tt.args.value)
			}
		})
	}
}

func Test_setUserProperty(t *testing.T) {
	type args struct {
		u        *api.User
		property string
		value    interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Sets the name property (string)",
			args: args{
				u:        defaultUser,
				property: "Name",
				value:    "Oompa Loompa",
			},
			wantErr: false,
		},
		{
			name: "Sets the MonitorGroups property ([]string)",
			args: args{
				u:        defaultUser,
				property: "MonitorGroups",
				value:    []string{"a", "b"},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setUserProperty(tt.args.u, tt.args.property, tt.args.value)
			if tt.args.property == "Name" && tt.args.u.Name != tt.args.value {
				t.Errorf("setUserProperty() name = %s, expected = %s", tt.args.u.Name, tt.args.value)
			}
			if tt.args.property == "MonitorGroups" && !reflect.DeepEqual(tt.args.u.MonitorGroups, tt.args.value) {
				t.Errorf("setUserProperty() MonitorGroups = %s, expected = %v", tt.args.u.MonitorGroups, tt.args.value)
			}
		})
	}
}
