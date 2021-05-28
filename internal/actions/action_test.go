// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: MIT

package actions_test

import (
	"testing"

	"github.com/vmware-tanzu-labs/yaml-overlay-tool/internal/actions"
	"gopkg.in/yaml.v3"
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
			name: "test Combine String",
			a:    actions.Combine,
			want: "combine",
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
		value *yaml.Node
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
				value: &yaml.Node{
					Kind:  yaml.ScalarNode,
					Value: "merge",
				},
			},
			wantErr: false,
		},
		{
			name: "Unmarshal Replace",
			a:    actions.Invalid,
			args: args{
				value: &yaml.Node{
					Kind:  yaml.ScalarNode,
					Value: "replace",
				},
			},
			wantErr: false,
		},
		{
			name: "Unmarshal Delete",
			a:    actions.Invalid,
			args: args{
				value: &yaml.Node{
					Kind:  yaml.ScalarNode,
					Value: "delete",
				},
			},
			wantErr: false,
		},
		{
			name: "Unmarshal Combine",
			a:    actions.Invalid,
			args: args{
				value: &yaml.Node{
					Kind:  yaml.ScalarNode,
					Value: "combine",
				},
			},
			wantErr: false,
		},
		{
			name: "Unmarshal Invalid Action",
			a:    actions.Invalid,
			args: args{
				value: &yaml.Node{
					Kind:  yaml.ScalarNode,
					Value: "potato",
				},
			},
			wantErr: true,
		},
		{
			name: "Unmarshal no Action",
			a:    actions.Invalid,
			args: args{
				value: &yaml.Node{
					Kind:  yaml.ScalarNode,
					Value: "",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		test := tt

		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			if err := test.a.UnmarshalYAML(test.args.value); (err != nil) != test.wantErr {
				t.Errorf("Action.UnmarshalYAML() error = %v, wantErr %v", err, test.wantErr)
			}
		})
	}
}
