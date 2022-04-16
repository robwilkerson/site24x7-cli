package user

import (
	"encoding/json"
	"errors"
	"reflect"
	"site24x7/api"
	"strings"
	"testing"

	"github.com/spf13/pflag"
)

func Test_list(t *testing.T) {
	mockAPIResponse := []byte(`[
		{"email_address": "foo@bar.com"},
		{"email_address": "humpty@dumpty.com"},
		{"email_address": "alice@wonderland.com"}
	]`)
	mockUserList := []api.User{
		{EmailAddress: "foo@bar.com"},
		{EmailAddress: "humpty@dumpty.com"},
		{EmailAddress: "alice@wonderland.com"},
	}

	tests := []struct {
		name       string
		apiListFn  func() (json.RawMessage, error)
		want       []api.User
		wantErr    bool
		wantErrMsg string
	}{
		{
			name: "Handles an API error",
			apiListFn: func() (json.RawMessage, error) {
				return nil, errors.New("testing!")
			},
			want:       nil,
			wantErr:    true,
			wantErrMsg: "testing!",
		},
		{
			name: "Returns a list of users",
			apiListFn: func() (json.RawMessage, error) {
				return mockAPIResponse, nil
			},
			want:    mockUserList,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		apiUserList = tt.apiListFn
		t.Run(tt.name, func(t *testing.T) {
			got, err := list()
			if (err != nil) != tt.wantErr {
				t.Errorf("list() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && !strings.Contains(err.Error(), tt.wantErrMsg) {
				t.Errorf("list() error = %v, wantErrMsg \"%s\"", err, tt.wantErrMsg)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("list() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_findByEmail(t *testing.T) {
	type args struct {
		email string
	}

	mockUserList := []api.User{
		{EmailAddress: "foo@bar.com"},
		{EmailAddress: "humpty@dumpty.com"},
		{EmailAddress: "alice@wonderland.com"},
	}

	tests := []struct {
		name       string
		args       args
		listFn     func() ([]api.User, error)
		want       *api.User
		wantErr    bool
		wantErrMsg string
	}{
		{
			name: "Handles an error from the list function",
			args: args{
				email: "aqua@man.com",
			},
			listFn: func() ([]api.User, error) {
				return nil, errors.New("testing!")
			},
			want:       nil,
			wantErr:    true,
			wantErrMsg: "testing!",
		},
		{
			name: "Throws a user not found error",
			args: args{
				email: "super@man.com",
			},
			listFn: func() ([]api.User, error) {
				return mockUserList, nil
			},
			want:       nil,
			wantErr:    true,
			wantErrMsg: "not found",
		},
		{
			name: "Returns a user with the given email addrss",
			args: args{
				email: "humpty@dumpty.com",
			},
			listFn: func() ([]api.User, error) {
				return mockUserList, nil
			},
			want:    &api.User{EmailAddress: "humpty@dumpty.com"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		list = tt.listFn
		t.Run(tt.name, func(t *testing.T) {
			got, err := findByEmail(tt.args.email)
			if (err != nil) != tt.wantErr {
				t.Errorf("findByEmail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && !strings.Contains(err.Error(), tt.wantErrMsg) {
				t.Errorf("findByEmail() error = %v, wantErrMsg \"%s\"", err, tt.wantErrMsg)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("findByEmail() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_get(t *testing.T) {
	type args struct {
		id    string
		email string
	}

	var mockUser api.User
	mockAPIResponse := []byte(`{
		"user_id": "1001001SOS",
		"email_address": "oompa@loompa.com"
	}`)
	json.Unmarshal(mockAPIResponse, &mockUser)

	tests := []struct {
		name       string
		args       args
		findFn     func(email string) (*api.User, error)
		apiGetFn   func(userID string) (json.RawMessage, error)
		want       *api.User
		wantErr    bool
		wantErrMsg string
	}{
		{
			name: "By email: Handles a finder error",
			args: args{
				id:    "",
				email: "foo@bar.com",
			},
			findFn: func(email string) (*api.User, error) {
				return nil, errors.New("testing!")
			},
			want:       nil,
			wantErr:    true,
			wantErrMsg: "testing!",
		},
		{
			name: "By email: Returns formatted json if found",
			args: args{
				id:    "",
				email: "foo@bar.com",
			},
			findFn: func(email string) (*api.User, error) {
				return &mockUser, nil
			},
			want:    &mockUser,
			wantErr: false,
		},
		{
			name: "By ID: Handles an API error",
			args: args{
				id:    "1001001SOS",
				email: "",
			},
			apiGetFn: func(userID string) (json.RawMessage, error) {
				return nil, errors.New("testing!")
			},
			want:       nil,
			wantErr:    true,
			wantErrMsg: "testing!",
		},
		{
			name: "By ID: Returns formatted JSON",
			args: args{
				id:    "1001001SOS",
				email: "",
			},
			apiGetFn: func(userID string) (json.RawMessage, error) {
				return mockAPIResponse, nil
			},
			want:    &mockUser,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		findByEmail = tt.findFn
		apiUserGet = tt.apiGetFn
		t.Run(tt.name, func(t *testing.T) {
			got, err := get(tt.args.id, tt.args.email)
			if (err != nil) != tt.wantErr {
				t.Errorf("get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && !strings.Contains(err.Error(), tt.wantErrMsg) {
				t.Errorf("get() error = %v, wantErrMsg \"%s\"", err, tt.wantErrMsg)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGet(t *testing.T) {
	type args struct {
		fs *pflag.FlagSet
	}

	e := pflag.NewFlagSet("emailOnlyTestFlags", pflag.ExitOnError)
	e.String("id", "", "")
	e.String("email", "dumb@luck.com", "")

	i := pflag.NewFlagSet("idOnlyTestFlags", pflag.ExitOnError)
	i.String("id", "1001001SOS", "")
	i.String("email", "", "")

	mockUser := &api.User{ID: "1001001SOS", EmailAddress: "oompa@loompa.com"}
	mockUserJSON, _ := json.MarshalIndent(mockUser, "", "    ")

	tests := []struct {
		name       string
		args       args
		getFn      func(id string, email string) (*api.User, error)
		want       []byte
		wantErr    bool
		wantErrMsg string
	}{
		{
			name: "Handles an error from the getter",
			args: args{
				fs: e,
			},
			getFn: func(id string, email string) (*api.User, error) {
				return nil, errors.New("testing!")
			},
			want:       nil,
			wantErr:    true,
			wantErrMsg: "testing!",
		},
		{
			name: "Returns formatted json if found",
			args: args{
				fs: e,
			},
			getFn: func(id string, email string) (*api.User, error) {
				return mockUser, nil
			},
			want:    mockUserJSON,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		get = tt.getFn
		t.Run(tt.name, func(t *testing.T) {
			got, err := Get(tt.args.fs)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && !strings.Contains(err.Error(), tt.wantErrMsg) {
				t.Errorf("Get() error = %v, wantErrMsg \"%s\"", err, tt.wantErrMsg)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestList(t *testing.T) {
	mockUserList := []api.User{
		{EmailAddress: "foo@bar.com"},
		{EmailAddress: "humpty@dumpty.com"},
		{EmailAddress: "alice@wonderland.com"},
	}
	mockUserJSON, _ := json.MarshalIndent(mockUserList, "", "    ")

	tests := []struct {
		name       string
		listFn     func() ([]api.User, error)
		want       []byte
		wantErr    bool
		wantErrMsg string
	}{
		{
			name: "Handles an error from the list function",
			listFn: func() ([]api.User, error) {
				return nil, errors.New("testing!")
			},
			want:       nil,
			wantErr:    true,
			wantErrMsg: "testing!",
		},
		{
			name: "Returns a list of users",
			listFn: func() ([]api.User, error) {
				return mockUserList, nil
			},
			want:    mockUserJSON,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		list = tt.listFn
		t.Run(tt.name, func(t *testing.T) {
			got, err := List()
			if (err != nil) != tt.wantErr {
				t.Errorf("List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && !strings.Contains(err.Error(), tt.wantErrMsg) {
				t.Errorf("List() error = %v, wantErrMsg \"%s\"", err, tt.wantErrMsg)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("List() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCreate(t *testing.T) {
	type args struct {
		email string
		fs    *pflag.FlagSet
	}

	mockAPIResponse := []byte(`{
		"is_edit_allowed":true,
		"is_client_portal_user":false,
		"is_account_contact":false,
		"is_invited":true,
		"subscribe_newsletter":false,
		"alert_settings":{
			"email_format":1,
			"anomaly":[1],
			"dont_alert_on_days":[],
			"critical":[1],
			"applogs":[1],
			"trouble":[1],
			"up":[1],
			"alerting_period": {
				"start_time":"00:00",
				"end_time":"00:00"
			},
			"down":[1]
		},
		"display_name":"Unnamed User",
		"twitter_settings":{},
		"is_contact":false,
		"selection_type":0,
		"user_role":0,
		"email_address":"faewtar234tgf2e@robwilkerson.org",
		"user_id":"366600000005462003",
		"mobile_settings":{},
		"image_present":false,
		"job_title":0,
		"notify_medium":[1],
		"user_groups":[]
	}`)
	var mockUser api.User
	json.Unmarshal(mockAPIResponse, &mockUser)
	mockUserJSON, _ := json.MarshalIndent(mockUser, "", "    ")

	var GetWriterFlagsWithInvalidProperty = func() *pflag.FlagSet {
		fs := GetWriterFlags()

		fs.String("random-invalid-property", "DUMMY", "Should be ignored")

		return fs
	}

	tests := []struct {
		name        string
		args        args
		apiCreateFn func(u *api.User) (json.RawMessage, error)
		want        []byte
		wantErr     bool
		wantErrMsg  string
	}{
		{
			name: "Handles an API error",
			args: args{
				email: "oompa@loompa.com",
				fs:    GetWriterFlags(),
			},
			apiCreateFn: func(u *api.User) (json.RawMessage, error) {
				return nil, errors.New("testing!")
			},
			want:       nil,
			wantErr:    true,
			wantErrMsg: "testing!",
		},
		{
			name: "Creates a default user",
			args: args{
				email: "boo@berry.com",
				fs:    GetWriterFlags(),
			},
			apiCreateFn: func(u *api.User) (json.RawMessage, error) {
				return mockAPIResponse, nil
			},
			want:    mockUserJSON,
			wantErr: false,
		},
		{
			name: "Ignores an invalid property",
			args: args{
				email: "boo@berry.com",
				fs:    GetWriterFlagsWithInvalidProperty(),
			},
			apiCreateFn: func(u *api.User) (json.RawMessage, error) {
				return mockAPIResponse, nil
			},
			want:    mockUserJSON,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		apiUserCreate = tt.apiCreateFn
		t.Run(tt.name, func(t *testing.T) {
			got, err := Create(tt.args.email, tt.args.fs)
			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && !strings.Contains(err.Error(), tt.wantErrMsg) {
				t.Errorf("Create() error = %v, wantErrMsg \"%s\"", err, tt.wantErrMsg)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Create() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUpdate(t *testing.T) {
	type args struct {
		fs *pflag.FlagSet
	}

	fs := GetAccessorFlags()
	fs.AddFlagSet(GetWriterFlags())

	mockUser := &api.User{ID: "1001001SOS", EmailAddress: "dizzy@dean.com"}
	mockUserUpdated := &api.User{
		ID:            "1001001SOS",
		Name:          "Dizzy Dean",
		EmailAddress:  "dizzy@dean.com",
		MonitorGroups: []string{"a", "b", "c"},
		AlertSettings: api.AlertSettings{
			EmailFormat: 0,
			SkipDays:    []int{6, 0},
			AlertingPeriod: api.AlertingPeriod{
				StartTime: "22:00",
			},
		},
		MobileSettings: api.MobileSettings{
			PhoneNumber: "3209957992",
		},
	}
	mockUserUpdatedPrettyJSON, _ := json.MarshalIndent(mockUserUpdated, "", "    ")

	tests := []struct {
		name            string
		args            args
		before          func()
		getFn           func(id string, email string) (*api.User, error)
		apiUserUpdateFn func(u *api.User) (json.RawMessage, error)
		want            []byte
		wantErr         bool
		wantErrMsg      string
	}{
		{
			name: "Handles an error from the get function",
			args: args{
				fs: fs,
			},
			before: func() {
				// noop
			},
			getFn: func(id string, email string) (*api.User, error) {
				return nil, errors.New("testing!")
			},
			want:       nil,
			wantErr:    true,
			wantErrMsg: "testing!",
		},
		{
			name: "Handles an API error",
			args: args{
				fs: fs,
			},
			before: func() {
				// noop
			},
			getFn: func(id string, email string) (*api.User, error) {
				return mockUser, nil
			},
			apiUserUpdateFn: func(u *api.User) (json.RawMessage, error) {
				return nil, errors.New("testing!")
			},
			want:       nil,
			wantErr:    true,
			wantErrMsg: "testing!",
		},
		{
			name: "Updates an existing user",
			args: args{
				fs: fs,
			},
			before: func() {
				// The .Set() function flips the .Changed flag
				fs.Set("name", "Dizzy Dean")
				fs.Set("monitor-groups", "a")
				fs.Set("monitor-groups", "b")
				fs.Set("monitor-groups", "c")
				fs.Set("alert-email-format", "0")
				fs.Set("alert-start-time", "22:00")
				fs.Set("alert-skip-days", "6")
				fs.Set("alert-skip-days", "0")
				fs.Set("mobile-phone-number", "3209957992")
				// this flag should be ignored
				fs.Set("non-eu-alert-consent", "true")
			},
			getFn: func(id string, email string) (*api.User, error) {
				return mockUser, nil
			},
			apiUserUpdateFn: func(u *api.User) (json.RawMessage, error) {
				// return what was sent
				j, _ := json.Marshal(u)

				return j, nil
			},
			want:    mockUserUpdatedPrettyJSON,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		get = tt.getFn
		apiUserUpdate = tt.apiUserUpdateFn
		t.Run(tt.name, func(t *testing.T) {
			tt.before()
			got, err := Update(tt.args.fs)
			if (err != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && !strings.Contains(err.Error(), tt.wantErrMsg) {
				t.Errorf("Update() error = %v, wantErrMsg \"%s\"", err, tt.wantErrMsg)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Update() = %v, want %v", string(got), string(tt.want))
			}
		})
	}
}

func TestDelete(t *testing.T) {
	type args struct {
		fs *pflag.FlagSet
	}

	fs := GetAccessorFlags()

	tests := []struct {
		name        string
		args        args
		getFn       func(id string, email string) (*api.User, error)
		apiDeleteFn func(id string) error
		wantErr     bool
		wantErrMsg  string
	}{
		{
			name: "Handles an error from the getter",
			args: args{
				fs: fs,
			},
			getFn: func(id string, email string) (*api.User, error) {
				return nil, errors.New("testing!")
			},
			apiDeleteFn: func(id string) error {
				// noop
				return nil
			},
			wantErr:    true,
			wantErrMsg: "testing!",
		},
		{
			name: "Handles an API error",
			args: args{
				fs: fs,
			},
			getFn: func(id string, email string) (*api.User, error) {
				return &api.User{ID: "1001001SOS"}, nil
			},
			apiDeleteFn: func(id string) error {
				return errors.New("testing!")
			},
			wantErr:    true,
			wantErrMsg: "testing!",
		},
		{
			name: "Returns successfully",
			args: args{
				fs: fs,
			},
			getFn: func(id string, email string) (*api.User, error) {
				return &api.User{ID: "1001001SOS"}, nil
			},
			apiDeleteFn: func(id string) error {
				return nil
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		get = tt.getFn
		apiUserDelete = tt.apiDeleteFn
		t.Run(tt.name, func(t *testing.T) {
			err := Delete(tt.args.fs)
			if (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr && !strings.Contains(err.Error(), tt.wantErrMsg) {
				t.Errorf("Delete() error = %v, wantErrMsg \"%s\"", err, tt.wantErrMsg)
			}
		})
	}
}
