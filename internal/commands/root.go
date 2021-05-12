// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: MIT

package commands

import (
	"fmt"
	"os"

	"github.com/op/go-logging"
	"github.com/spf13/cobra"
	"github.com/vmware-tanzu-labs/yaml-overlay-tool/internal/lib"
)

var ErrMissingRequired = fmt.Errorf("missing required arguments")

func New() *cobra.Command {
	// rootCmd represents the base command when called without any subcommands.
	rootCmd := &cobra.Command{
		Use:                        "yot",
		Aliases:                    []string{},
		SuggestFor:                 []string{},
		Short:                      YotShort,
		Long:                       YotLong,
		Example:                    "",
		ValidArgs:                  []string{},
		ValidArgsFunction:          nil,
		Args:                       nil,
		ArgAliases:                 []string{},
		BashCompletionFunction:     "",
		Deprecated:                 "",
		Annotations:                map[string]string{},
		Version:                    Version,
		PersistentPreRun:           SetupLogging,
		PersistentPreRunE:          nil,
		PreRun:                     nil,
		PreRunE:                    nil,
		Run:                        nil,
		RunE:                       Execute,
		PostRun:                    nil,
		PostRunE:                   nil,
		PersistentPostRun:          nil,
		PersistentPostRunE:         nil,
		FParseErrWhitelist:         cobra.FParseErrWhitelist{},
		TraverseChildren:           false,
		Hidden:                     false,
		SilenceErrors:              false,
		SilenceUsage:               false,
		DisableFlagParsing:         false,
		DisableAutoGenTag:          false,
		DisableFlagsInUseLine:      false,
		DisableSuggestions:         false,
		SuggestionsMinimumDistance: 0,
	}

	initializeGlobalFlags(rootCmd)

	rootCmd.AddCommand(addCompletionCommand())

	return rootCmd
}

func SetupLogging(cmd *cobra.Command, args []string) {
	cmd.SetOut(cmd.OutOrStdout())

	format := logging.MustStringFormatter(
		`%{color}%{time:15:04:05} [%{level}]%{color:reset} %{message}`,
	)
	backend := logging.AddModuleLevel(
		logging.NewBackendFormatter(logging.NewLogBackend(os.Stderr, "", 0), format))

	if options.Verbose {
		backend.SetLevel(logging.DEBUG, "")
	} else {
		backend.SetLevel(logging.ERROR, "")
	}

	logging.SetBackend(backend)
}

func Execute(cmd *cobra.Command, args []string) error {
	if err := lib.Execute(&options); err != nil {
		cmd.SilenceUsage = true

		return fmt.Errorf("%w", err)
	}

	return nil
}
