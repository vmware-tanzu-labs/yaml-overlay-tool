// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: MIT

package commands

import (
	"os"

	"github.com/spf13/cobra"
)

func (r *Root) CompletionCommand() *cobra.Command {
	return &cobra.Command{
		Use:                   completionUse,
		Short:                 completionShort,
		Long:                  completionLong,
		DisableFlagsInUseLine: true,
		ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
		Args:                  cobra.ExactValidArgs(1),
		Run:                   r.Completions,
	}
}

func (r *Root) Completions(cmd *cobra.Command, args []string) {
	if len(args) > 1 {
		panic("Too Many Args")
	}

	switch args[0] {
	case "bash":
		if err := cmd.Root().GenBashCompletion(os.Stdout); err != nil {
			r.Log.Fatal(err)
		}
	case "zsh":
		if err := cmd.Root().GenZshCompletion(os.Stdout); err != nil {
			r.Log.Fatal(err)
		}
	case "fish":
		if err := cmd.Root().GenFishCompletion(os.Stdout, true); err != nil {
			r.Log.Fatal(err)
		}
	case "powershell":
		if err := cmd.Root().GenPowerShellCompletionWithDesc(os.Stdout); err != nil {
			r.Log.Fatal(err)
		}
	}
}
