// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: MIT

package commands

import (
	"fmt"
	"os"

	"github.com/op/go-logging"
	"github.com/spf13/cobra"
	"github.com/vmware-tanzu-labs/yaml-overlay-tool/internal/actions"
	"github.com/vmware-tanzu-labs/yaml-overlay-tool/internal/instructions"
)

// ErrMissingRequired occurs when a required flag is not passed.
var ErrMissingRequired = fmt.Errorf("missing required arguments")

type Command struct {
	Log     *logging.Logger
	Options *instructions.Config
}

type RootCommand Command

func New() *RootCommand {
	return &RootCommand{
		Log: logging.MustGetLogger("cmd"),
		Options: &instructions.Config{
			LogLevel: logging.ERROR,
			Styles:   actions.Styles{actions.NormalStyle},
		},
	}
}

func (rc RootCommand) Command(version string) *cobra.Command {
	// rootCmd represents the base command when called without any subcommands.
	rootCmd := &cobra.Command{
		Use:                        "yot",
		Aliases:                    []string{},
		SuggestFor:                 []string{},
		Short:                      YotShort,
		Long:                       YotLong,
		Example:                    HelpUsageExample,
		ValidArgs:                  []string{},
		Args:                       nil,
		ArgAliases:                 []string{},
		BashCompletionFunction:     "",
		Deprecated:                 "",
		Hidden:                     false,
		Annotations:                map[string]string{},
		Version:                    version,
		PersistentPreRun:           rc.SetupLogging,
		PersistentPreRunE:          nil,
		PreRun:                     nil,
		PreRunE:                    nil,
		Run:                        rc.Execute,
		RunE:                       nil,
		PostRun:                    nil,
		PostRunE:                   nil,
		PersistentPostRun:          nil,
		PersistentPostRunE:         nil,
		SilenceErrors:              false,
		SilenceUsage:               false,
		DisableFlagParsing:         false,
		DisableAutoGenTag:          false,
		DisableFlagsInUseLine:      false,
		DisableSuggestions:         false,
		SuggestionsMinimumDistance: 0,
		TraverseChildren:           false,
		FParseErrWhitelist:         cobra.FParseErrWhitelist{},
	}

	rc.initializeGlobalFlags(rootCmd)

	rootCmd.AddCommand(CompletionCommand(rc).Command())

	return rootCmd
}

func (rc *RootCommand) SetupLogging(cmd *cobra.Command, args []string) {
	cmd.SetOut(cmd.OutOrStdout())

	format := logging.MustStringFormatter(
		`%{color}%{time:15:04:05} [%{level}]%{color:reset} %{message}`,
	)
	backend := logging.AddModuleLevel(
		logging.NewBackendFormatter(logging.NewLogBackend(os.Stderr, "", 0), format),
	)

	backend.SetLevel(rc.Options.LogLevel, "")
	logging.SetBackend(backend)
}

func (rc *RootCommand) Execute(cmd *cobra.Command, args []string) {
	if err := instructions.Execute(rc.Options); err != nil {
		cmd.SilenceUsage = true

		rc.Log.Error(fmt.Errorf("%w", err))
		os.Exit(1)
	}
}
