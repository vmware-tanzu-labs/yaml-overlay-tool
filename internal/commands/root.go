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

type Root struct {
	Log     *logging.Logger
	Options *instructions.Config
	Command *cobra.Command
}

func New() *Root {
	rc := &Root{
		Log: logging.MustGetLogger("cmd"),
		Options: &instructions.Config{
			LogLevel: logging.ERROR,
			Styles:   actions.Styles{actions.NormalStyle},
		},
	}

	rc.Command = rc.NewCommand()
	rc.AddFlags()
	rc.AddCommands()

	return rc
}

func (r Root) NewCommand() *cobra.Command {
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
		PersistentPreRun:           r.SetupLogging,
		PersistentPreRunE:          nil,
		PreRun:                     nil,
		PreRunE:                    nil,
		Run:                        r.Execute,
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

	return rootCmd
}

func (r *Root) AddFlags() {
	r.initializeGlobalFlags()
	r.initializeTemplateFlags()
}

func (r *Root) AddCommands() {
	r.Command.AddCommand(r.CompletionCommand())
}

func (r *Root) SetupLogging(cmd *cobra.Command, args []string) {
	cmd.SetOut(cmd.OutOrStdout())

	format := logging.MustStringFormatter(
		`%{color}%{time:15:04:05} [%{level}]%{color:reset} %{message}`,
	)
	backend := logging.AddModuleLevel(
		logging.NewBackendFormatter(logging.NewLogBackend(os.Stderr, "", 0), format),
	)

	backend.SetLevel(r.Options.LogLevel, "")
	logging.SetBackend(backend)
}

func (r *Root) Execute(cmd *cobra.Command, args []string) {
	if err := instructions.Execute(r.Options); err != nil {
		cmd.SilenceUsage = true

		r.Log.Error(fmt.Errorf("%w", err))
		os.Exit(1)
	}
}

func (r *Root) Run() {
	if err := r.Command.Execute(); err != nil {
		r.Log.Error(err)
	}
}
