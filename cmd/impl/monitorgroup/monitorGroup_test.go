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
