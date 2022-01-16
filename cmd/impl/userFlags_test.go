//
// User types that are - or can be - shared across multiple user<Action> files
//
package impl

import (
	"reflect"
	"testing"

	"github.com/spf13/pflag"
)

func TestPropertyMapper(t *testing.T) {
	fs := pflag.NewFlagSet("testFlags", pflag.ExitOnError)
	type args struct {
		f    *pflag.FlagSet
		name string
	}
	tests := []struct {
		name string
		args args
		want pflag.NormalizedName
	}{
		// default handler
		{
			name: "Normalizes a single word flag name",
			args: args{
				f:    fs,
				name: "name",
			},
			want: "Name",
		},
		{
			name: "Normalizes a multi-word flag name",
			args: args{
				f:    fs,
				name: "multi-word-flag",
			},
			want: "MultiWordFlag",
		},
		// special case handlers - internal nomenclature translation
		{
			name: "Normalizes the notify-by flag",
			args: args{
				f:    fs,
				name: "notify-by",
			},
			want: "NotificationMethods",
		},
		// special case handlers - nested user properties
		{
			name: "Normalizes the alert-methods-down flag",
			args: args{
				f:    fs,
				name: "alert-methods-down",
			},
			want: "AlertDownNotificationMethods",
		},
		{
			name: "Normalizes the alert-methods-trouble flag",
			args: args{
				f:    fs,
				name: "alert-methods-trouble",
			},
			want: "AlertTroubleNotificationMethods",
		},
		{
			name: "Normalizes the alert-methods-up flag",
			args: args{
				f:    fs,
				name: "alert-methods-up",
			},
			want: "AlertUpNotificationMethods",
		},
		{
			name: "Normalizes the alert-methods-applogs flag",
			args: args{
				f:    fs,
				name: "alert-methods-applogs",
			},
			want: "AlertAppLogsNotificationMethods",
		},
		{
			name: "Normalizes the alert-methods-anomaly flag",
			args: args{
				f:    fs,
				name: "alert-methods-anomaly",
			},
			want: "AlertAnomalyNotificationMethods",
		},
		// special case handlers - acronyms
		{
			name: "Normalizes the statusiq-role flag",
			args: args{
				f:    fs,
				name: "statusiq-role",
			},
			want: "StatusIQRole",
		},
		{
			name: "Normalizes the mobile-sms-provider-id flag",
			args: args{
				f:    fs,
				name: "mobile-sms-provider-id",
			},
			want: "MobileSMSProviderID",
		},
		{
			name: "Normalizes the mobile-call-provider-id flag",
			args: args{
				f:    fs,
				name: "mobile-call-provider-id",
			},
			want: "MobileCallProviderID",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := PropertyMapper(tt.args.f, tt.args.name); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PropertyMapper() = %v, want %v", got, tt.want)
			}
		})
	}
}
