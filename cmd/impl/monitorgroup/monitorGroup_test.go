package monitorgroup

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
	type args struct {
		withSubgroups bool
	}

	mockAPIResponse := []byte(`[
		{"display_name": "Test 1"},
		{"display_name": "Test 2"},
		{"display_name": "Test 3"}
	]`)
	mockList := []api.MonitorGroup{
		{Name: "Test 1"},
		{Name: "Test 2"},
		{Name: "Test 3"},
	}

	tests := []struct {
		name       string
		args       args
		apiListFn  func(withSubgroups bool) (json.RawMessage, error)
		want       []api.MonitorGroup
		wantErr    bool
		wantErrMsg string
	}{
		{
			name: "Handles an API error",
			apiListFn: func(withSubgroups bool) (json.RawMessage, error) {
				return nil, errors.New("testing!")
			},
			want:       nil,
			wantErr:    true,
			wantErrMsg: "testing!",
		},
		{
			name: "Returns a list of monitor groups",
			apiListFn: func(withSubgroups bool) (json.RawMessage, error) {
				return mockAPIResponse, nil
			},
			want:    mockList,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		apiMonitorGroupList = tt.apiListFn
		t.Run(tt.name, func(t *testing.T) {
			got, err := list(tt.args.withSubgroups)
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

func TestList(t *testing.T) {
	type args struct {
		fs *pflag.FlagSet
	}

	mockList := []api.MonitorGroup{
		{Name: "Test 1"},
		{Name: "Test 2"},
		{Name: "Test 3"},
	}
	mockJSON, _ := json.MarshalIndent(mockList, "", "    ")

	tests := []struct {
		name       string
		args       args
		listFn     func(withSubgroups bool) ([]api.MonitorGroup, error)
		want       []byte
		wantErr    bool
		wantErrMsg string
	}{
		{
			name: "Handles an error from the list function",
			args: args{
				fs: pflag.NewFlagSet("testing", pflag.PanicOnError),
			},
			listFn: func(withSubgroups bool) ([]api.MonitorGroup, error) {
				return nil, errors.New("testing!")
			},
			want:       nil,
			wantErr:    true,
			wantErrMsg: "testing!",
		},
		{
			name: "Returns a list of users",
			args: args{
				fs: pflag.NewFlagSet("testing", pflag.PanicOnError),
			},
			listFn: func(withSubgroups bool) ([]api.MonitorGroup, error) {
				return mockList, nil
			},
			want:    mockJSON,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		list = tt.listFn
		t.Run(tt.name, func(t *testing.T) {
			got, err := List(tt.args.fs)
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
		name string
		fs   *pflag.FlagSet
	}

	mockAPIResponse := []byte(`{
		"display_name": "Test Group",
		"description": "Just for testing",
		"monitors": [],
		"dependency_resource_ids":[]
	}`)
	var mockMonitorGroup api.MonitorGroup
	json.Unmarshal(mockAPIResponse, &mockMonitorGroup)
	mockMonitorGroupJSON, _ := json.MarshalIndent(mockMonitorGroup, "", "    ")

	tests := []struct {
		name        string
		args        args
		apiCreateFn func(u *api.MonitorGroup) (json.RawMessage, error)
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
			apiCreateFn: func(u *api.MonitorGroup) (json.RawMessage, error) {
				return nil, errors.New("testing!")
			},
			want:       nil,
			wantErr:    true,
			wantErrMsg: "testing!",
		},
		{
			name: "Creates a monitor group",
			args: args{
				name: "Test Group",
				fs:   GetWriterFlags(),
			},
			apiCreateFn: func(u *api.MonitorGroup) (json.RawMessage, error) {
				return mockAPIResponse, nil
			},
			want:    mockMonitorGroupJSON,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		apiMonitorGroupCreate = tt.apiCreateFn
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
		fs *pflag.FlagSet
	}

	mockFlagSet := pflag.NewFlagSet("Testing", pflag.PanicOnError)

	var mockMonitorGroup api.MonitorGroup
	mockAPIResponse := []byte(`{
		"group_id": "1001001SOS",
		"display_name": "TESTING",
		"description": "Test Group",
		"health_threshold_count": 1,
		"group_type": 1
	}`)
	json.Unmarshal(mockAPIResponse, &mockMonitorGroup)
	mockJSON, _ := json.MarshalIndent(mockMonitorGroup, "", "    ")

	tests := []struct {
		name       string
		args       args
		apiGetFn   func(id string) (json.RawMessage, error)
		want       []byte
		wantErr    bool
		wantErrMsg string
	}{
		{
			name: "Handles an API error",
			args: args{
				id: "1001001SOS",
				fs: mockFlagSet,
			},
			apiGetFn: func(id string) (json.RawMessage, error) {
				return nil, errors.New("testing!")
			},
			want:       nil,
			wantErr:    true,
			wantErrMsg: "testing!",
		},
		{
			name: "Returns formatted json",
			args: args{
				id: "1001001SOS",
				fs: mockFlagSet,
			},
			apiGetFn: func(id string) (json.RawMessage, error) {
				return mockAPIResponse, nil
			},
			want:    mockJSON,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		apiMonitorGroupGet = tt.apiGetFn
		t.Run(tt.name, func(t *testing.T) {
			got, err := Get(tt.args.id, tt.args.fs)
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
