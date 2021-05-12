// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: MIT

package commands

import (
	"testing"

	"github.com/spf13/cobra"
)

func Test_initializeGlobalFlags(t *testing.T) {
	type args struct {
		rootCmd *cobra.Command
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test initialization",
			args: args{
				&cobra.Command{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.rootCmd.SetArgs([]string{"-v"})
			initializeGlobalFlags(tt.args.rootCmd)
		})
	}
}
