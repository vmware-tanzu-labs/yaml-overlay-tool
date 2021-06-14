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
				Use:              commands.YotUse,
				Aliases:          []string{},
				SuggestFor:       []string{},
				Short:            commands.YotShort,
				Long:             commands.YotLong,
				Version:          "unstable",
				PersistentPreRun: commands.New().SetupLogging,
				RunE:             commands.New().Execute,
			},
		},
	}

	for _, tt := range tests {
		test := tt
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			got := commands.New()
			if got.Command.Use != test.want.Use {
				t.Errorf("New() = %v, want %v", got.Command.Use, test.want.Use)
			}
			if got.Command.Short != test.want.Short {
				t.Errorf("New() = %v, want %v", got.Command.Short, test.want.Short)
			}
			if got.Command.Long != test.want.Long {
				t.Errorf("New() = %v, want %v", got.Command.Long, test.want.Long)
			}
			if got.Command.Version != test.want.Version {
				t.Errorf("New() = %v, want %v", got.Command.Version, test.want.Version)
			}
			if reflect.DeepEqual(got.Execute, test.want.Execute) {
				t.Errorf("New() = %v, want %v", got.Command.Version, test.want.Version)
			}
		})
	}
}
