// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: MIT

package commands

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/op/go-logging"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/vmware-tanzu-labs/yaml-overlay-tool/internal/actions"
	"github.com/vmware-tanzu-labs/yaml-overlay-tool/internal/instructions"
)

// ErrMissingRequired occurs when a required flag is not passed.
var ErrMissingRequired = fmt.Errorf("missing required arguments")

type Root struct {
	configFile string
	Log        *logging.Logger
	Options    *instructions.Config
	Command    *cobra.Command
}

func New() *Root {
	rc := &Root{
		Log: logging.MustGetLogger("cmd"),
		Options: &instructions.Config{
			LogLevel: logging.ERROR,
			Styles:   actions.Styles{actions.NormalStyle},
		},
	}

	cobra.OnInitialize(rc.initConfig)

	rc.Command = rc.NewCommand()
	rc.AddFlags()
	rc.AddCommands()

	return rc
}

func (r *Root) initConfig() {
	if r.configFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(r.configFile)
	} else {
		viper.SetConfigName(".yot")
		viper.SetConfigType("yaml")
		viper.AddConfigPath("/etc/appname/")  // path to look for the config file in
		viper.AddConfigPath("$HOME/.appname") // call multiple times to add many search paths
		viper.AddConfigPath(".")
	}

	viper.SetEnvPrefix("yot")

	replacer := strings.NewReplacer("-", "_")

	viper.SetEnvKeyReplacer(replacer)

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if ok := errors.As(err, &viper.ConfigFileNotFoundError{}); !ok {
			r.Log.Fatal(err)
		}

		return
	}

	r.Log.Debugf("Using config file: %s", viper.ConfigFileUsed())
}

func (r Root) NewCommand() *cobra.Command {
	// rootCmd represents the base command when called without any subcommands.
	rootCmd := &cobra.Command{
		Use:                        "yot",
		Aliases:                    []string{},
		SuggestFor:                 []string{},
		Short:                      YotShort,
		Long:                       YotLong,
		Example:                    helpUsageExample,
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
		RunE:                       r.Execute,
		PostRun:                    nil,
		PostRunE:                   nil,
		PersistentPostRun:          nil,
		PersistentPostRunE:         nil,
		SilenceErrors:              false,
		SilenceUsage:               false,
		DisableFlagParsing:         false,
		DisableAutoGenTag:          true,
		DisableFlagsInUseLine:      false,
		DisableSuggestions:         false,
		SuggestionsMinimumDistance: 0,
		TraverseChildren:           false,
		FParseErrWhitelist:         cobra.FParseErrWhitelist{},
	}

	return rootCmd
}

func (r *Root) AddFlags() {
	r.initializeCommonFlags()
	r.initializeInstructionFlags()
	r.initializeOutputFlags()
	r.initializeFormatFlags()
	r.initializeTemplateFlags()
	r.initializeStdInFlags()
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

	logLevel, err := logging.LogLevel((viper.GetString("log-level")))
	if err != nil {
		panic(err)
	}

	backend.SetLevel(logLevel, "")

	logging.SetBackend(backend)
}

func (r *Root) Execute(cmd *cobra.Command, args []string) error {
	if err := instructions.Execute(r.Options); err != nil {
		cmd.SilenceUsage = true

		return fmt.Errorf("%w", err)
	}

	return nil
}

func (r *Root) Run() {
	if err := r.Command.Execute(); err != nil {
		r.Log.Error(err)
		os.Exit(1)
	}
}
