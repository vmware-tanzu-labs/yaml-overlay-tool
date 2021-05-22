// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: MIT

package commands_test

import (
	"fmt"
	"testing"

	"github.com/spf13/cobra"
	"github.com/vmware-tanzu-labs/yaml-overlay-tool/internal/commands"
)

func TestCompletions(t *testing.T) {
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
			name: "Without any args",
			args: args{
				cmd:  commands.CompletionCommand{}.Command(),
				args: nil,
			},
			wantErr: true,
		},
		{
			name: "With bash arg",
			args: args{
				cmd:  commands.CompletionCommand{}.Command(),
				args: []string{"bash"},
			},
			wantErr: false,
		},
		{
			name: "Without Zsh arg",
			args: args{
				cmd:  commands.CompletionCommand{}.Command(),
				args: []string{"zsh"},
			},
			wantErr: false,
		},
		{
			name: "With fish arg",
			args: args{
				cmd:  commands.CompletionCommand{}.Command(),
				args: []string{"fish"},
			},
			wantErr: false,
		},
		{
			name: "With powershell arg",
			args: args{
				cmd:  commands.CompletionCommand{}.Command(),
				args: []string{"powershell"},
			},
			wantErr: false,
		},
		{
			name: "Without too many args",
			args: args{
				cmd:  commands.CompletionCommand{}.Command(),
				args: []string{"bash", "zsh"},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		test := tt
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			defer func() {
				if err := recover(); err != nil {
					newErr := fmt.Errorf("%w", err)

					if !test.wantErr && newErr != nil {
						t.Errorf("Completions() error = %v, wantErr %v", err, test.wantErr)
					}
				}
			}()
			commands.CompletionCommand{}.Completions(test.args.cmd, test.args.args)
			if test.wantErr {
				t.Errorf("Completions() error = %v, wantErr %v", nil, test.wantErr)
			}
		})
	}
}
