package user

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/spf13/pflag"
)

func Test_lookup(t *testing.T) {
	type args struct {
		keys   []int
		lookup map[int]string
	}
	mockMap := map[int]string{
		0: "a",
		1: "b",
		2: "c",
		3: "d",
		4: "e",
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			name: "All values exist in the map",
			args: args{
				keys:   []int{0, 1, 2, 3, 4},
				lookup: mockMap,
			},
			want: []int{0, 1, 2, 3, 4},
		},
		{
			name: "No values exist in the map",
			args: args{
				keys:   []int{5, 6, 7, 8},
				lookup: mockMap,
			},
			want: nil,
		},
		{
			name: "Some values exist in the map",
			args: args{
				keys:   []int{8, 2, 3, 9, 12, 14},
				lookup: mockMap,
			},
			want: []int{2, 3},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := lookup(tt.args.keys, tt.args.lookup); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("lookup() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_validateAccessors(t *testing.T) {
	z := pflag.NewFlagSet("emptyTestFlags", pflag.ExitOnError)
	z.String("id", "", "")
	z.String("email", "", "")

	i := pflag.NewFlagSet("idOnlyTestFlags", pflag.ExitOnError)
	i.String("id", "20392059", "")
	i.String("email", "", "")

	e := pflag.NewFlagSet("emailOnlyTestFlags", pflag.ExitOnError)
	e.String("id", "", "")
	e.String("email", "foo@bar.com", "")

	b := pflag.NewFlagSet("bothTestFlags", pflag.ExitOnError)
	b.String("id", "2i9e02932", "")
	b.String("email", "foo@bar.com", "")

	type args struct {
		fs *pflag.FlagSet
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "No flag values sent",
			args: args{
				fs: z,
			},
			wantErr: true,
		},
		{
			name: "ID value sent",
			args: args{
				fs: i,
			},
			wantErr: false,
		},
		{
			name: "Email value sent",
			args: args{
				fs: e,
			},
			wantErr: false,
		},
		{
			name: "ID & email values sent",
			args: args{
				fs: b,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateAccessors(tt.args.fs); (err != nil) != tt.wantErr {
				t.Errorf("validateAccessors() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func getDefaultFlags() *pflag.FlagSet {
	d := pflag.NewFlagSet("defaultTestFlags", pflag.ExitOnError)
	d.Int("role", 0, "")
	d.Int("job-title", 0, "")
	d.IntSlice("notify-by", []int{1}, "")
	d.Int("resource-type", 0, "")
	d.Int("alert-email-format", 1, "")
	d.IntSlice("alert-skip-days", []int{}, "")
	d.IntSlice("alert-methods-down", []int{1}, "")
	d.IntSlice("alert-methods-trouble", []int{1}, "")
	d.IntSlice("alert-methods-up", []int{1}, "")
	d.IntSlice("alert-methods-applogs", []int{1}, "")
	d.IntSlice("alert-methods-anomaly", []int{1}, "")
	d.Int("statusiq-role", 0, "")
	d.Int("cloudspend-role", 0, "")

	return d
}

func Test_validateWriters(t *testing.T) {
	type args struct {
		fs *pflag.FlagSet
	}
	tests := []struct {
		name      string
		args      args
		before    func(*pflag.FlagSet)
		wantPanic bool
	}{
		{
			name: "Default values are valid",
			args: args{
				fs: getDefaultFlags(),
			},
			before: func(fs *pflag.FlagSet) {
				// noop
			},
			wantPanic: false,
		},
		{
			name: "Invalid role",
			args: args{
				fs: getDefaultFlags(),
			},
			before: func(fs *pflag.FlagSet) {
				fs.Set("role", "-1")
			},
			wantPanic: true,
		},
		{
			name: "Invalid job-title",
			args: args{
				fs: getDefaultFlags(),
			},
			before: func(fs *pflag.FlagSet) {
				fs.Set("job-title", "-1")
			},
			wantPanic: true,
		},
		{
			name: "Invalid notify-by",
			args: args{
				fs: getDefaultFlags(),
			},
			before: func(fs *pflag.FlagSet) {
				fs.Set("notify-by", "500")
				fs.Set("notify-by", "-1")
			},
			wantPanic: true,
		},
		{
			name: "Invalid resource-type",
			args: args{
				fs: getDefaultFlags(),
			},
			before: func(fs *pflag.FlagSet) {
				fs.Set("resource-type", "-1")
			},
			wantPanic: true,
		},
		{
			name: "Invalid alert-email-format",
			args: args{
				fs: getDefaultFlags(),
			},
			before: func(fs *pflag.FlagSet) {
				fs.Set("alert-email-format", "-1")
			},
			wantPanic: true,
		},
		{
			name: "Invalid alert-skip-days (too many)",
			args: args{
				fs: getDefaultFlags(),
			},
			before: func(fs *pflag.FlagSet) {
				fs.Set("alert-skip-days", "0")
				fs.Set("alert-skip-days", "1")
				fs.Set("alert-skip-days", "2")
				fs.Set("alert-skip-days", "3")
				fs.Set("alert-skip-days", "4")
				fs.Set("alert-skip-days", "5")
				fs.Set("alert-skip-days", "6")
				fs.Set("alert-skip-days", "7")
			},
			wantPanic: true,
		},
		{
			name: "Invalid alert-skip-days (value > 6)",
			args: args{
				fs: getDefaultFlags(),
			},
			before: func(fs *pflag.FlagSet) {
				fs.Set("alert-skip-days", "0")
				fs.Set("alert-skip-days", "8")
			},
			wantPanic: true,
		},
		{
			name: "Invalid alert-skip-days (value < 0)",
			args: args{
				fs: getDefaultFlags(),
			},
			before: func(fs *pflag.FlagSet) {
				fs.Set("alert-skip-days", "-1")
				fs.Set("alert-skip-days", "4")
			},
			wantPanic: true,
		},
		{
			name: "Invalid alert-methods-down",
			args: args{
				fs: getDefaultFlags(),
			},
			before: func(fs *pflag.FlagSet) {
				fs.Set("alert-methods-down", "-1")
			},
			wantPanic: true,
		},
		{
			name: "Invalid alert-methods-trouble",
			args: args{
				fs: getDefaultFlags(),
			},
			before: func(fs *pflag.FlagSet) {
				fs.Set("alert-methods-trouble", "-1")
			},
			wantPanic: true,
		},
		{
			name: "Invalid alert-methods-up",
			args: args{
				fs: getDefaultFlags(),
			},
			before: func(fs *pflag.FlagSet) {
				fs.Set("alert-methods-up", "-1")
			},
			wantPanic: true,
		},
		{
			name: "Invalid alert-methods-applogs",
			args: args{
				fs: getDefaultFlags(),
			},
			before: func(fs *pflag.FlagSet) {
				fs.Set("alert-methods-applogs", "-1")
			},
			wantPanic: true,
		},
		{
			name: "Invalid alert-methods-anomaly",
			args: args{
				fs: getDefaultFlags(),
			},
			before: func(fs *pflag.FlagSet) {
				fs.Set("alert-methods-anomaly", "-1")
			},
			wantPanic: true,
		},
		{
			name: "Invalid statusiq-role",
			args: args{
				fs: getDefaultFlags(),
			},
			before: func(fs *pflag.FlagSet) {
				fs.Set("statusiq-role", "-1")
			},
			wantPanic: true,
		},
		{
			name: "Invalid cloudspend-role",
			args: args{
				fs: getDefaultFlags(),
			},
			before: func(fs *pflag.FlagSet) {
				fs.Set("cloudspend-role", "-1")
			},
			wantPanic: true,
		},
	}
	for _, tt := range tests {
		fmt.Println(strings.ToUpper(tt.name))
		t.Run(tt.name, func(t *testing.T) {
			// No need to check whether `recover()` is nil. Just turn off the panic.
			defer func() { recover() }()

			// call each test's before() function to update a particular flag
			// value before attempting to validate
			tt.before(tt.args.fs)

			validateWriters(tt.args.fs)

			// Never reaches here if `validateWriters` panics. No panic should
			// happen for valid writers so we explicitly tell each test whether
			// we expect it to panic.
			if tt.wantPanic {
				t.Errorf("expected a panic")
			}
		})
	}
}

func Test_normalizeName(t *testing.T) {
	type args struct {
		f *pflag.Flag
	}

	d := pflag.NewFlagSet("defaultTestFlags", pflag.ExitOnError)
	// default cases
	d.Int("simpleflag", 0, "")
	d.Int("multi-word-flag", 0, "")
	// special cases
	d.IntSlice("notify-by", []int{1}, "")
	d.Int("alert-start-time", 0, "")
	d.Int("alert-end-time", 0, "")
	// nested properties
	d.Int("alert-email-format", 1, "")
	d.IntSlice("alert-skip-days", []int{}, "")
	d.IntSlice("alert-methods-down", []int{1}, "")
	d.IntSlice("alert-methods-trouble", []int{1}, "")
	d.IntSlice("alert-methods-up", []int{1}, "")
	d.IntSlice("alert-methods-applogs", []int{1}, "")
	d.IntSlice("alert-methods-anomaly", []int{1}, "")
	// acronym handling
	d.Int("statusiq-role", 0, "")
	d.Int("mobile-sms-provider-id", 0, "")
	d.Int("mobile-call-provider-id", 0, "")

	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Handles a single word, simple flag",
			args: args{
				f: d.Lookup("simpleflag"),
			},
			want: "Simpleflag",
		},
		{
			name: "Handles a multi-word, simple flag",
			args: args{
				f: d.Lookup("multi-word-flag"),
			},
			want: "MultiWordFlag",
		},
		{
			name: "Handles notify-by",
			args: args{
				f: d.Lookup("notify-by"),
			},
			want: "NotificationMethods",
		},
		{
			name: "Handles alert-start-time",
			args: args{
				f: d.Lookup("alert-start-time"),
			},
			want: "AlertingPeriodStartTime",
		},
		{
			name: "Handles alert-end-time",
			args: args{
				f: d.Lookup("alert-end-time"),
			},
			want: "AlertingPeriodEndTime",
		},
		{
			name: "Handles alert-methods-down",
			args: args{
				f: d.Lookup("alert-methods-down"),
			},
			want: "AlertDownNotificationMethods",
		},
		{
			name: "Handles alert-methods-trouble",
			args: args{
				f: d.Lookup("alert-methods-trouble"),
			},
			want: "AlertTroubleNotificationMethods",
		},
		{
			name: "Handles alert-methods-up",
			args: args{
				f: d.Lookup("alert-methods-up"),
			},
			want: "AlertUpNotificationMethods",
		},
		{
			name: "Handles alert-methods-applogs",
			args: args{
				f: d.Lookup("alert-methods-applogs"),
			},
			want: "AlertAppLogsNotificationMethods",
		},
		{
			name: "Handles alert-methods-anomaly",
			args: args{
				f: d.Lookup("alert-methods-anomaly"),
			},
			want: "AlertAnomalyNotificationMethods",
		},
		{
			name: "Handles statusiq-role",
			args: args{
				f: d.Lookup("statusiq-role"),
			},
			want: "StatusIQRole",
		},
		{
			name: "Handles mobile-sms-provider-id",
			args: args{
				f: d.Lookup("mobile-sms-provider-id"),
			},
			want: "MobileSMSProviderID",
		},
		{
			name: "Handles mobile-call-provider-id",
			args: args{
				f: d.Lookup("mobile-call-provider-id"),
			},
			want: "MobileCallProviderID",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := normalizeName(tt.args.f); got != tt.want {
				t.Errorf("normalizeName() = %v, want %v", got, tt.want)
			}
		})
	}
}
