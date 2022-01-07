//
// Implementation and supporting functions for the `user create` subcommand.
//
package impl

import (
	"fmt"
	"reflect"
	"site24x7/api"
	"strings"
	"testing"
)

func Test_lookupIds(t *testing.T) {
	type args struct {
		list   []int
		lookup map[int]string
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			name: "Key exists",
			args: args{
				[]int{1},
				map[int]string{
					0: "a",
					1: "b",
				},
			},
			want: []int{1},
		},
		{
			name: "Key doesn't exist",
			args: args{
				[]int{5},
				map[int]string{0: "a", 1: "b", 2: "c"},
			},
			want: nil,
		},
		{
			name: "Multiple keys exist",
			args: args{
				[]int{1, 3, 11, 49, 10},
				map[int]string{
					0:  "a",
					1:  "b",
					2:  "c",
					3:  "d",
					10: "x",
					11: "7",
				},
			},
			want: []int{1, 3, 11, 10},
		},
		{
			name: "No keys exist",
			args: args{
				[]int{33, 18, 41, 9},
				map[int]string{
					0:  "a",
					1:  "b",
					2:  "c",
					3:  "d",
					10: "x",
					11: "y",
				},
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := lookupIds(tt.args.list, tt.args.lookup); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("lookupIds() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_UserWriterFlags_validate(t *testing.T) {
	type fields struct {
		role                int
		notifyMethod        []int
		statusIQRole        int
		cloudSpendRole      int
		alertEmailFormat    int
		alertSkipDays       []int
		alertMethodsDown    []int
		alertMethodsTrouble []int
		alertMethodsUp      []int
		alertMethodsAppLogs []int
		alertMethodsAnomaly []int
		jobTitle            int
		resourceType        int
	}
	// defaultFlags sets the default value for only those flags that are
	// validated; the "name" flag, for example, is not validated and is not
	// represented in this stuct.
	defaultFlags := fields{
		role:                0,
		notifyMethod:        []int{1},
		statusIQRole:        0,
		cloudSpendRole:      0,
		alertEmailFormat:    1,
		alertSkipDays:       []int{},
		alertMethodsDown:    []int{1},
		alertMethodsTrouble: []int{1},
		alertMethodsUp:      []int{1},
		alertMethodsAppLogs: []int{1},
		alertMethodsAnomaly: []int{1},
		jobTitle:            0,
		resourceType:        0,
	}

	// various invalid flag states
	invalidRoleFlag := defaultFlags
	invalidRoleFlag.role = 45
	invalidNotificationMethods := defaultFlags
	invalidNotificationMethods.notifyMethod = []int{100, 1001}
	invalidStatusIQRole := defaultFlags
	invalidStatusIQRole.statusIQRole = 23902
	invalidCloudSpendRole := defaultFlags
	invalidCloudSpendRole.cloudSpendRole = 329
	invalidAlertEmailFormat := defaultFlags
	invalidAlertEmailFormat.alertEmailFormat = 5
	invalidAlertSkipDaysLongWeek := defaultFlags
	invalidAlertSkipDaysLongWeek.alertSkipDays = []int{9, 4, 8, 3, 21, 47, 12, 1}
	invalidAlertSkipDaysLessThanZero := defaultFlags
	invalidAlertSkipDaysLessThanZero.alertSkipDays = []int{4, 1, -1}
	invalidAlertSkipDaysGreaterThanSix := defaultFlags
	invalidAlertSkipDaysGreaterThanSix.alertSkipDays = []int{1, 7}
	invalidAlertDownNotification := defaultFlags
	invalidAlertDownNotification.alertMethodsDown = []int{-1, 500}
	invalidAlertTroubleNotification := defaultFlags
	invalidAlertTroubleNotification.alertMethodsTrouble = []int{-1, 500}
	invalidAlertUpNotification := defaultFlags
	invalidAlertUpNotification.alertMethodsUp = []int{500}
	invalidAlertAppLogsNotification := defaultFlags
	invalidAlertAppLogsNotification.alertMethodsAppLogs = []int{239}
	invalidAlertAnomalyNotification := defaultFlags
	invalidAlertAnomalyNotification.alertMethodsAnomaly = []int{19}
	invalidJobTitle := defaultFlags
	invalidJobTitle.jobTitle = 10
	invalidResourceType := defaultFlags
	invalidResourceType.resourceType = 10

	tests := []struct {
		name       string
		fields     fields
		wantErr    bool
		wantErrMsg string
	}{
		{
			name:       "Default flag values are valid",
			fields:     defaultFlags,
			wantErr:    false,
			wantErrMsg: "",
		},
		{
			name:       "An unsupported role throws an error",
			fields:     invalidRoleFlag,
			wantErr:    true,
			wantErrMsg: "ERROR: Invalid role",
		},
		{
			name:       "No valid notification methods were passed",
			fields:     invalidNotificationMethods,
			wantErr:    true,
			wantErrMsg: "ERROR: Invalid notification method(s)",
		},
		{
			name:       "An unempty, unsupported status IQ role is invalid",
			fields:     invalidStatusIQRole,
			wantErr:    true,
			wantErrMsg: "ERROR: Invalid status IQ role",
		},
		{
			name:       "An unempty, unsupported cloudspend role is invalid",
			fields:     invalidCloudSpendRole,
			wantErr:    true,
			wantErrMsg: "ERROR: Invalid cloudspend role",
		},
		{
			name:       "An unsupported email format throws an error",
			fields:     invalidAlertEmailFormat,
			wantErr:    true,
			wantErrMsg: "ERROR: Invalid email format",
		},
		{
			name:       "More skip days than days in a week throws an error",
			fields:     invalidAlertSkipDaysLongWeek,
			wantErr:    true,
			wantErrMsg: "ERROR: There are 7 days in a week",
		},
		{
			name:       "Any skip days value < 0 throws an error",
			fields:     invalidAlertSkipDaysLessThanZero,
			wantErr:    true,
			wantErrMsg: "ERROR: Invalid skip days identified",
		},
		{
			name:       "Any skip days value > 6 throws an error",
			fields:     invalidAlertSkipDaysGreaterThanSix,
			wantErr:    true,
			wantErrMsg: "ERROR: Invalid skip days identified",
		},
		{
			name:       "No valid DOWN notification methods were passed",
			fields:     invalidAlertDownNotification,
			wantErr:    true,
			wantErrMsg: "ERROR: Invalid DOWN alert notification method(s)",
		},
		{
			name:       "No valid TROUBLE notification methods were passed",
			fields:     invalidAlertTroubleNotification,
			wantErr:    true,
			wantErrMsg: "ERROR: Invalid TROUBLE alert notification method(s)",
		},
		{
			name:       "No valid UP notification methods were passed",
			fields:     invalidAlertUpNotification,
			wantErr:    true,
			wantErrMsg: "ERROR: Invalid UP alert notification method(s)",
		},
		{
			name:       "No valid APPLOGS notification methods were passed",
			fields:     invalidAlertAppLogsNotification,
			wantErr:    true,
			wantErrMsg: "ERROR: Invalid APPLOGS alert notification method(s)",
		},
		{
			name:       "No valid ANOMALY notification methods were passed",
			fields:     invalidAlertAnomalyNotification,
			wantErr:    true,
			wantErrMsg: "ERROR: Invalid ANOMALY alert notification method(s)",
		},
		{
			name:       "An invalid job title throws an error",
			fields:     invalidJobTitle,
			wantErr:    true,
			wantErrMsg: "ERROR: Invalid job title",
		},
		{
			name:       "An invalid resource type thrown an error",
			fields:     invalidResourceType,
			wantErr:    true,
			wantErrMsg: "ERROR: Invalid resource type",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := UserWriterFlags{
				Role:                tt.fields.role,
				NotifyMethod:        tt.fields.notifyMethod,
				StatusIQRole:        tt.fields.statusIQRole,
				CloudSpendRole:      tt.fields.cloudSpendRole,
				AlertEmailFormat:    tt.fields.alertEmailFormat,
				AlertSkipDays:       tt.fields.alertSkipDays,
				AlertMethodsDown:    tt.fields.alertMethodsDown,
				AlertMethodsTrouble: tt.fields.alertMethodsTrouble,
				AlertMethodsUp:      tt.fields.alertMethodsUp,
				AlertMethodsAppLogs: tt.fields.alertMethodsAppLogs,
				AlertMethodsAnomaly: tt.fields.alertMethodsAnomaly,
				JobTitle:            tt.fields.jobTitle,
				ResourceType:        tt.fields.resourceType,
			}
			err := f.validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("UserWriterFlags.validate() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil && !strings.HasPrefix(err.Error(), tt.wantErrMsg) {
				t.Errorf("UserWriterFlags.validate() error msg = %s, wantErrMsg = %s", err.Error(), tt.wantErrMsg)
			}
		})
	}
}

func Test_userCreate(t *testing.T) {
	type args struct {
		f       UserWriterFlags
		u       *api.User
		creator func() error
	}
	defaultFlags := UserWriterFlags{
		Name:                "Unnamed User",
		Role:                0,
		NotifyMethod:        []int{1},
		StatusIQRole:        0,
		CloudSpendRole:      0,
		AlertEmailFormat:    1,
		AlertSkipDays:       []int{},
		AlertStartTime:      "00:00",
		AlertEndTime:        "00:00",
		AlertMethodsDown:    []int{1},
		AlertMethodsTrouble: []int{1},
		AlertMethodsUp:      []int{1},
		AlertMethodsAppLogs: []int{1},
		AlertMethodsAnomaly: []int{1},
		JobTitle:            0,
		ResourceType:        0,
	}
	mockInvalidFlagset := defaultFlags
	mockInvalidFlagset.Role = -1
	mockFlagsetWithStatusIQRole := defaultFlags
	mockFlagsetWithStatusIQRole.StatusIQRole = 25
	mockFlagsetWithCloudspendRole := defaultFlags
	mockFlagsetWithCloudspendRole.CloudSpendRole = 12

	// A mock user as it would exist entering the creator function
	defaultAlertMethods := []int{1}
	defaultAlertSettings := api.UserAlertSettings{
		EmailFormat: 1,
		SkipDays:    []int{},
		AlertingPeriod: api.UserAlertingPeriod{
			StartTime: "00:00",
			EndTime:   "00:00",
		},
		DownAlertMethods:    defaultAlertMethods,
		TroubleAlertMethods: defaultAlertMethods,
		UpAlertMethods:      defaultAlertMethods,
		AppLogsAlertMethods: defaultAlertMethods,
		AnomalyAlertMethods: defaultAlertMethods,
	}
	mockEntryUser := &api.User{EmailAddress: "super@man.com"}
	// A mock user as it would get updated if no flag values were explicitly
	// passed
	mockHydratedDefaultUser := &api.User{
		Name:               "Unnamed User",
		EmailAddress:       "super@man.com",
		Role:               0,
		NotificationMethod: []int{1},
		AlertSettings:      defaultAlertSettings,
		JobTitle:           0,
		MobileSettings: map[string]interface{}{
			"call_provider_id": 0,
			"country_code":     "",
			"mobile_number":    "",
			"sms_provider_id":  0,
		},
		ResourceType: 0,
	}
	// A hydrated mock user where a non-empty, non-default statusiq-role flag
	// value has been passed
	mockHydratedUserWithCustomStatusIQRole := &api.User{
		Name:               "Unnamed User",
		EmailAddress:       "super@man.com",
		Role:               0,
		NotificationMethod: []int{1},
		AlertSettings:      defaultAlertSettings,
		JobTitle:           0,
		StatusIQRole:       25,
		MobileSettings: map[string]interface{}{
			"call_provider_id": 0,
			"country_code":     "",
			"mobile_number":    "",
			"sms_provider_id":  0,
		},
	}
	// A hydrated mock user where a non-empty, non-default cloudspend-role flag
	// value has been passed
	mockHydratedUserWithCustomCloudspendRole := &api.User{
		Name:               "Unnamed User",
		EmailAddress:       "super@man.com",
		Role:               0,
		NotificationMethod: []int{1},

		AlertSettings: defaultAlertSettings,
		JobTitle:      0,
		MobileSettings: map[string]interface{}{
			"call_provider_id": 0,
			"country_code":     "",
			"mobile_number":    "",
			"sms_provider_id":  0,
		},
		CloudspendRole: 12,
	}
	tests := []struct {
		name       string
		args       args
		wantUser   *api.User
		wantErr    bool
		wantErrMsg string
	}{
		{
			name: "Rethrows a flag validation error",
			args: args{
				f: mockInvalidFlagset,
				u: mockEntryUser,
				creator: func() error {
					return nil
				},
			},
			wantUser:   mockEntryUser,
			wantErr:    true,
			wantErrMsg: "ERROR: Invalid role",
		},
		{
			name: "Hydrates the passed user object with default values",
			args: args{
				f: defaultFlags,
				u: mockEntryUser,
				creator: func() error {
					return nil
				},
			},
			wantUser: mockHydratedDefaultUser,
			wantErr:  false,
		},
		{
			name: "Rethrows an error returned from the creator function",
			args: args{
				f: defaultFlags,
				u: mockEntryUser,
				creator: func() error {
					return fmt.Errorf("Whoops!")
				},
			},
			wantUser:   mockHydratedDefaultUser,
			wantErr:    true,
			wantErrMsg: "Whoops!",
		},
		{
			name: "Sets the StatusIQ role if a valid, non-default value is passed",
			args: args{
				f: mockFlagsetWithStatusIQRole,
				u: mockEntryUser,
				creator: func() error {
					return nil
				},
			},
			wantUser: mockHydratedUserWithCustomStatusIQRole,
			wantErr:  false,
		},
		{
			name: "Sets the Cloudspend role if a valid, non-default value is passed",
			args: args{
				f: mockFlagsetWithCloudspendRole,
				u: mockEntryUser,
				creator: func() error {
					return nil
				},
			},
			wantUser: mockHydratedUserWithCustomCloudspendRole,
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := UserCreate(tt.args.f, tt.args.u, tt.args.creator)
			if (err != nil) != tt.wantErr {
				t.Errorf("userCreate() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil && !strings.HasPrefix(err.Error(), tt.wantErrMsg) {
				t.Errorf("UserWriterFlags.validate() error msg = %s, wantErrMsg = %s", err.Error(), tt.wantErrMsg)
			}
			if tt.wantUser != nil && !reflect.DeepEqual(tt.args.u, tt.wantUser) {
				t.Errorf("userCreate() = %+v, want %+v", tt.args.u, tt.wantUser)
			}
		})
	}
}
