// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: MIT

package commands_test

import (
	"reflect"
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
				PersistentPreRun: commands.New().SetupLogging,
				Run:              commands.New().Execute,
			},
		},
	}

	for _, tt := range tests {
		test := tt
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			got := commands.New().Command()
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
			if reflect.DeepEqual(got.Execute, test.want.Execute) {
				t.Errorf("New() = %v, want %v", got.Version, test.want.Version)
			}
		})
	}
}
