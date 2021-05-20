// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: MIT

package actions_test

import (
	"reflect"
	"testing"

	"github.com/vmware-tanzu-labs/yaml-overlay-tool/internal/actions"
)

func TestAction_String(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		a    actions.Action
		want string
	}{
		{
			name: "Test Merge String",
			a:    actions.Merge,
			want: "merge",
		},
		{
			name: "Test Delete String",
			a:    actions.Delete,
			want: "delete",
		},
		{
			name: "Test Replace String",
			a:    actions.Replace,
			want: "replace",
		},
		{
			name: "test Math String",
			a:    actions.Math,
			want: "math",
		},
		{
			name: "Test invalid String",
			a:    actions.Invalid,
			want: "",
		},
	}
	for _, tt := range tests {
		test := tt
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			if got := test.a.String(); got != test.want {
				t.Errorf("Action.String() = %v, want %v", got, test.want)
			}
		})
	}
}

func TestAction_UnmarshalYAML(t *testing.T) {
	t.Parallel()

	type args struct {
		unmarshal func(interface{}) error
	}

	tests := []struct {
		name    string
		a       actions.Action
		args    args
		wantErr bool
	}{
		{
			name: "Unmarshal Merge",
			a:    actions.Invalid,
			args: args{
				unmarshal: testUnmarshal("merge"),
			},
			wantErr: false,
		},
		{
			name: "Unmarshal Replace",
			a:    actions.Invalid,
			args: args{
				unmarshal: testUnmarshal("replace"),
			},
			wantErr: false,
		},
		{
			name: "Unmarshal Delete",
			a:    actions.Invalid,
			args: args{
				unmarshal: testUnmarshal("delete"),
			},
			wantErr: false,
		},
		{
			name: "Unmarshal Math",
			a:    actions.Invalid,
			args: args{
				unmarshal: testUnmarshal("math"),
			},
			wantErr: false,
		},
		{
			name: "Unmarshal Invalid Action",
			a:    actions.Invalid,
			args: args{
				unmarshal: testUnmarshal("potato"),
			},
			wantErr: true,
		},
		{
			name: "Unmarshal no Action",
			a:    actions.Invalid,
			args: args{
				unmarshal: testUnmarshal(""),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		test := tt

		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			if err := test.a.UnmarshalYAML(test.args.unmarshal); (err != nil) != test.wantErr {
				t.Errorf("Action.UnmarshalYAML() error = %v, wantErr %v", err, test.wantErr)
			}
		})
	}
}

func testUnmarshal(str string) func(interface{}) error {
	return func(y interface{}) error {
		testString := str
		if s, ok := y.(*string); ok {
			*s = testString
		} else {
			return actions.ErrInvalidAction
		}

		return nil
	}
}

func TestAction_MarshalYAML(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		a       actions.Action
		want    interface{}
		wantErr bool
	}{
		{
			name:    "Marshal Merge",
			a:       actions.Merge,
			want:    "merge",
			wantErr: false,
		},
		{
			name:    "Unmarshal Replace",
			a:       actions.Replace,
			want:    "replace",
			wantErr: false,
		},
		{
			name:    "Unmarshal Delete",
			a:       actions.Delete,
			want:    "delete",
			wantErr: false,
		},
		{
			name:    "Unmarshal Math",
			a:       actions.Math,
			want:    "math",
			wantErr: false,
		},
		{
			name:    "Unmarshal Invalid Action",
			a:       actions.Invalid,
			want:    "",
			wantErr: false,
		},
		{
			name:    "Unmarshal no Action",
			a:       actions.Invalid,
			want:    "",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		test := tt

		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			got, err := test.a.MarshalYAML()
			if (err != nil) != test.wantErr {
				t.Errorf("Action.MarshalYAML() error = %v, wantErr %v", err, test.wantErr)

				return
			}
			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("Action.MarshalYAML() = %v, want %v", got, test.want)
			}
		})
	}
}

func TestOnMissingAction_String(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		a    actions.OnMissingAction
		want string
	}{
		{
			name: "Test Merge String",
			a:    actions.Ignore,
			want: "ignore",
		},
		{
			name: "Test Delete String",
			a:    actions.Inject,
			want: "inject",
		},
		{
			name: "Test invalid String",
			a:    actions.Invalid,
			want: "ignore",
		},
	}
	for _, tt := range tests {
		test := tt
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			if got := test.a.String(); got != test.want {
				t.Errorf("OnMissingAction.String() = %v, want %v", got, test.want)
			}
		})
	}
}

func TestOnMissingAction_UnmarshalYAML(t *testing.T) {
	t.Parallel()

	type args struct {
		unmarshal func(interface{}) error
	}

	tests := []struct {
		name    string
		a       actions.OnMissingAction
		args    args
		wantErr bool
	}{
		{
			name: "Unmarshal Merge",
			a:    actions.Ignore,
			args: args{
				unmarshal: testUnmarshal("ignore"),
			},
			wantErr: false,
		},
		{
			name: "Unmarshal Replace",
			a:    actions.Ignore,
			args: args{
				unmarshal: testUnmarshal("Inject"),
			},
			wantErr: false,
		},
		{
			name: "Unmarshal Invalid Action",
			a:    actions.Invalid,
			args: args{
				unmarshal: testUnmarshal("potato"),
			},
			wantErr: false,
		},
		{
			name: "Unmarshal no Action",
			a:    actions.Invalid,
			args: args{
				unmarshal: testUnmarshal(""),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		test := tt

		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			if err := test.a.UnmarshalYAML(test.args.unmarshal); (err != nil) != test.wantErr {
				t.Errorf("OnMissingAction.UnmarshalYAML() error = %v, wantErr %v", err, test.wantErr)
			}
		})
	}
}

func TestOnMissingAction_MarshalYAML(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		a       actions.OnMissingAction
		want    interface{}
		wantErr bool
	}{
		{
			name:    "Marshal Ignore",
			a:       actions.Ignore,
			want:    "ignore",
			wantErr: false,
		},
		{
			name:    "Unmarshal Inject",
			a:       actions.Inject,
			want:    "inject",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		test := tt

		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			got, err := test.a.MarshalYAML()
			if (err != nil) != test.wantErr {
				t.Errorf("Action.MarshalYAML() error = %v, wantErr %v", err, test.wantErr)

				return
			}
			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("Action.MarshalYAML() = %v, want %v", got, test.want)
			}
		})
	}
}
