package usergroup

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
		{"display_name": "Team 1"},
		{"display_name": "Team 2"},
		{"display_name": "Team 3"}
	]`)
	mockAPIBadJSON := []byte(`[
		{"display_name": "Team 1"},
		{"display_name": "Team 2"},
		{"display_name": "Team 3"},
	]`)
	mockList := []api.UserGroup{
		{Name: "Team 1"},
		{Name: "Team 2"},
		{Name: "Team 3"},
	}

	tests := []struct {
		name       string
		apiListFn  func() (json.RawMessage, error)
		want       []api.UserGroup
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
			name: "Handles a JSON parsing error",
			apiListFn: func() (json.RawMessage, error) {
				return mockAPIBadJSON, nil
			},
			want:       nil,
			wantErr:    true,
			wantErrMsg: "Unable to  parse response data",
		},
		{
			name: "Returns a list of user groups",
			apiListFn: func() (json.RawMessage, error) {
				return mockAPIResponse, nil
			},
			want:    mockList,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		apiUserGroupList = tt.apiListFn
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

func Test_get(t *testing.T) {
	type args struct {
		id string
	}

	var mockUserGroup api.UserGroup
	mockAPIResponse := []byte(`{
		"user_group_id": "1001001SOS",
		"display_name": "Team 1"
	}`)
	mockBadAPIResponse := []byte(`{
		"user_group_id": "1001001SOS",
		"display_name": "Team 1",
	}`)
	json.Unmarshal(mockAPIResponse, &mockUserGroup)

	tests := []struct {
		name       string
		args       args
		findFn     func(id string) (*api.UserGroup, error)
		apiGetFn   func(id string) (json.RawMessage, error)
		want       *api.UserGroup
		wantErr    bool
		wantErrMsg string
	}{
		{
			name: "Handles an API error",
			args: args{
				id: "1001001SOS",
			},
			apiGetFn: func(test string) (json.RawMessage, error) {
				return nil, errors.New("testing!")
			},
			want:       nil,
			wantErr:    true,
			wantErrMsg: "testing!",
		},
		{
			name: "Handles a malformed API response",
			args: args{
				id: "1001001SOS",
			},
			apiGetFn: func(test string) (json.RawMessage, error) {
				return mockBadAPIResponse, nil
			},
			want:       nil,
			wantErr:    true,
			wantErrMsg: "Unable to  parse response data",
		},
		{
			name: "Returns formatted JSON",
			args: args{
				id: "1001001SOS",
			},
			apiGetFn: func(userID string) (json.RawMessage, error) {
				return mockAPIResponse, nil
			},
			want:    &mockUserGroup,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		apiUserGroupGet = tt.apiGetFn
		t.Run(tt.name, func(t *testing.T) {
			got, err := get(tt.args.id)
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

func TestCreate(t *testing.T) {
	type args struct {
		name string
		fs   *pflag.FlagSet
	}

	mockAPIResponse := []byte(`{
		"display_name": "Team 1",
		"users": ["foo", "bar", "baz"]
	}`)
	mockBadAPIResponse := []byte(`{
		"display_name": "Team 1",
		"users": ["foo", "bar", "baz"],
	}`)
	var mockUserGroup api.UserGroup
	json.Unmarshal(mockAPIResponse, &mockUserGroup)
	mockMonitorGroupJSON, _ := json.MarshalIndent(mockUserGroup, "", "    ")

	tests := []struct {
		name        string
		args        args
		apiCreateFn func(u *api.UserGroup) (json.RawMessage, error)
		want        []byte
		wantErr     bool
		wantErrMsg  string
	}{
		{
			name: "Handles an API error",
			args: args{
				name: "Test Group",
				fs:   GetWriterFlags(),
			},
			apiCreateFn: func(u *api.UserGroup) (json.RawMessage, error) {
				return nil, errors.New("Tried to create a team without users")
			},
			want:       nil,
			wantErr:    true,
			wantErrMsg: "team without users",
		},
		{
			name: "Handles an error parsing the response",
			args: args{
				name: "Team 1",
				fs:   GetWriterFlags(),
			},
			apiCreateFn: func(u *api.UserGroup) (json.RawMessage, error) {
				return mockBadAPIResponse, nil
			},
			want:       nil,
			wantErr:    true,
			wantErrMsg: "Unable to  parse response data",
		},
		{
			name: "Creates a user group",
			args: args{
				name: "Test Group",
				fs:   GetWriterFlags(),
			},
			apiCreateFn: func(u *api.UserGroup) (json.RawMessage, error) {
				return mockAPIResponse, nil
			},
			want:    mockMonitorGroupJSON,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		apiUserGroupCreate = tt.apiCreateFn
		t.Run(tt.name, func(t *testing.T) {
			got, err := Create(tt.args.name, tt.args.fs)
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

func TestGet(t *testing.T) {
	type args struct {
		id string
	}

	mockUserGroup := &api.UserGroup{ID: "1001001SOS", Name: "Team 1"}
	mockUserGroupJSON, _ := json.MarshalIndent(mockUserGroup, "", "    ")

	tests := []struct {
		name       string
		args       args
		getFn      func(id string) (*api.UserGroup, error)
		want       []byte
		wantErr    bool
		wantErrMsg string
	}{
		{
			name: "Handles an error from the getter",
			args: args{
				id: "1001001SOS",
			},
			getFn: func(id string) (*api.UserGroup, error) {
				return nil, errors.New("testing!")
			},
			want:       nil,
			wantErr:    true,
			wantErrMsg: "testing!",
		},
		{
			name: "Returns formatted json if found",
			args: args{
				id: "1001001SOS",
			},
			getFn: func(id string) (*api.UserGroup, error) {
				return mockUserGroup, nil
			},
			want:    mockUserGroupJSON,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		get = tt.getFn
		t.Run(tt.name, func(t *testing.T) {
			got, err := Get(tt.args.id)
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
	mockList := []api.UserGroup{
		{Name: "Team 1"},
		{Name: "Team 2"},
		{Name: "Team 3"},
	}
	mockJSON, _ := json.MarshalIndent(mockList, "", "    ")

	tests := []struct {
		name       string
		listFn     func() ([]api.UserGroup, error)
		want       []byte
		wantErr    bool
		wantErrMsg string
	}{
		{
			name: "Handles an error from the list function",
			listFn: func() ([]api.UserGroup, error) {
				return nil, errors.New("testing!")
			},
			want:       nil,
			wantErr:    true,
			wantErrMsg: "testing!",
		},
		{
			name: "Returns a list of users",
			listFn: func() ([]api.UserGroup, error) {
				return mockList, nil
			},
			want:    mockJSON,
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
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("List() = %v, want %v", got, tt.want)
			}
		})
	}
}
