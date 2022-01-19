package user

import (
	"encoding/json"
	"errors"
	"reflect"
	"site24x7/api"
	"testing"

	"github.com/spf13/pflag"
)

func Test_setProperty(t *testing.T) {
	type args struct {
		v        interface{}
		property string
		value    interface{}
	}
	expectedUser := &api.User{
		Name: "Humpty Dumpty",
	}
	expectedAlertSettings := &api.UserAlertSettings{
		SkipDays: []int{0, 6},
	}
	expectedAlertingPeriodSettings := &api.UserAlertingPeriod{
		StartTime: "08:30",
	}
	tests := []struct {
		name string
		args args
		want interface{}
	}{
		{
			name: "Sets a top-level User property",
			args: args{
				v:        &api.User{},
				property: "Name",
				value:    "Humpty Dumpty",
			},
			want: expectedUser,
		},
		{
			name: "Sets an AlertSettings property",
			args: args{
				v:        &api.UserAlertSettings{},
				property: "SkipDays",
				value:    []int{0, 6},
			},
			want: expectedAlertSettings,
		},
		{
			name: "Sets an AlertingPeriod property",
			args: args{
				v:        &api.UserAlertingPeriod{},
				property: "StartTime",
				value:    "08:30",
			},
			want: expectedAlertingPeriodSettings,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setProperty(tt.args.v, tt.args.property, tt.args.value)

			if !reflect.DeepEqual(tt.args.v, tt.want) {
				t.Errorf("Read() = %v, want %v", tt.args.v, tt.want)
			}
		})
	}
}

func TestCreate(t *testing.T) {
	type args struct {
		fs      *pflag.FlagSet
		u       *api.User
		creator func() error
	}

	d := pflag.NewFlagSet("testFlags", pflag.ExitOnError)
	// user properties
	d.String("name", "Jack of Spades", "")
	d.IntSlice("notify-by", []int{1, 2}, "")
	d.StringSlice("monitor-groups", []string{"a", "b", "c"}, "")
	// alerting period properties
	d.String("alert-start-time", "08:30", "")
	d.String("alert-end-time", "09:00", "")
	// alert settings
	d.Int("alert-email-format", 1, "")
	d.IntSlice("alert-skip-days", []int{0, 6}, "")
	// mobile settings
	d.String("mobile-phone-number", "5553909980", "")
	// ignored flags
	d.Int("statusiq-role", -1, "")
	d.String("help", "", "")

	expectedUser := &api.User{
		Name:                "Jack of Spades",
		NotificationMethods: []int{1, 2},
		MonitorGroups:       []string{"a", "b", "c"},
		AlertSettings: api.UserAlertSettings{
			EmailFormat: 1,
			SkipDays:    []int{0, 6},
			AlertingPeriod: api.UserAlertingPeriod{
				StartTime: "08:30",
				EndTime:   "09:00",
			},
		},
		MobileSettings: api.UserMobileSettings{
			PhoneNumber: "5553909980",
		},
	}
	expectedJSON, _ := json.MarshalIndent(expectedUser, "", "    ")

	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "Handles an error thrown by the creator function",
			args: args{
				fs: d,
				u:  &api.User{},
				creator: func() error {
					return errors.New("testing!")
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Returns a marshalled user struct",
			args: args{
				fs: d,
				u:  &api.User{},
				creator: func() error {
					return nil
				},
			},
			want:    expectedJSON,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Create(tt.args.fs, tt.args.u, tt.args.creator)
			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Create() = %v, want %v", got, tt.want)
			}
		})
	}
}
