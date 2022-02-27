package usergroup

import (
	"encoding/json"
	"errors"
	"reflect"
	"site24x7/api"
	"strings"
	"testing"
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
