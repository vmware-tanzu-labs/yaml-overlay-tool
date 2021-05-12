// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: MIT

package commands_test

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/vmware-tanzu-labs/yaml-overlay-tool/internal/commands"
)

func TestNew(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		want *cobra.Command
	}{
		{
			name: "check initialization",
			want: &cobra.Command{
				Use:              "yot",
				Aliases:          []string{},
				SuggestFor:       []string{},
				Short:            commands.YotShort,
				Long:             commands.YotLong,
				Version:          commands.Version,
				PersistentPreRun: commands.SetupLogging,
				RunE:             commands.Execute,
			},
		},
	}

	for _, tt := range tests {
		test := tt
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			got := commands.New()
			if got.Use != test.want.Use {
				t.Errorf("New() = %v, want %v", got.Use, test.want.Use)
			}
			if got.Short != test.want.Short {
				t.Errorf("New() = %v, want %v", got.Short, test.want.Short)
			}
			if got.Long != test.want.Long {
				t.Errorf("New() = %v, want %v", got.Long, test.want.Long)
			}
			if got.Version != test.want.Version {
				t.Errorf("New() = %v, want %v", got.Version, test.want.Version)
			}
		})
	}
}

func TestSetupLogging(t *testing.T) {
	t.Parallel()

	type args struct {
		cmd  *cobra.Command
		args []string
	}

	tests := []struct {
		name string
		args args
	}{
		{
			name: "Test Initialization",
			args: args{
				cmd:  &cobra.Command{},
				args: nil,
			},
		},
	}

	for _, tt := range tests {
		test := tt

		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			commands.SetupLogging(test.args.cmd, test.args.args)
		})
	}
}

func TestExecute(t *testing.T) {
	t.Parallel()

	type args struct {
		cmd  *cobra.Command
		args []string
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "No args",
			args: args{
				cmd:  commands.New(),
				args: nil,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		test := tt
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			if err := commands.Execute(test.args.cmd, test.args.args); (err != nil) != test.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, test.wantErr)
			}
		})
	}
}
