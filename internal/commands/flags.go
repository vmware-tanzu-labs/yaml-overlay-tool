// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: MIT

package commands

import (
	"github.com/op/go-logging"
	"github.com/spf13/viper"
	"github.com/thediveo/enumflag"
	"github.com/vmware-tanzu-labs/yaml-overlay-tool/internal/actions"
)

func (r *Root) initializeCommonFlags() {
	r.Command.Flags().StringVar(&r.configFile, "config", r.configFile, "config file (default is $HOME/.yot)")

	if err := viper.BindPFlag("configFile", r.Command.Flags().Lookup("config")); err != nil {
		r.Log.Fatal(err)
	}

	logMap := map[logging.Level][]string{
		logging.CRITICAL: {"critical", "crit", "c"},
		logging.ERROR:    {"error", "err", "e"},
		logging.WARNING:  {"warning", "warn", "w"},
		logging.NOTICE:   {"notice", "note", "n"},
		logging.INFO:     {"info", "i"},
		logging.DEBUG:    {"debug", "d", "verbose", "v"},
	}

	r.Command.Flags().VarP(
		enumflag.New(&r.Options.LogLevel, "logLevel", logMap, enumflag.EnumCaseInsensitive),
		"log-level",
		"v",
		helpLogLevel,
	)

	r.Command.Flags().Lookup("log-level").NoOptDefVal = "debug"

	if err := viper.BindPFlag("logLevel", r.Command.Flags().Lookup("log-level")); err != nil {
		r.Log.Fatal(err)
	}
}

func (r *Root) initializeOutputFlags() {
	r.Command.Flags().StringVarP(
		&r.Options.OutputDir,
		"output-directory",
		"o",
		"./output",
		helpOutputDirectory,
	)

	if err := r.Command.MarkFlagDirname("output-directory"); err != nil {
		r.Log.Fatal(err)
	}

	r.Command.Flags().BoolVarP(
		&r.Options.StdOut,
		"stdout",
		"s",
		false,
		helpRenderStdOut,
	)

	if err := viper.BindPFlag("outputDirectory", r.Command.Flags().Lookup("output-directory")); err != nil {
		r.Log.Fatal(err)
	}

	if err := viper.BindPFlag("stdout", r.Command.Flags().Lookup("stdout")); err != nil {
		r.Log.Fatal(err)
	}

	viper.SetDefault("outputDirectory", "./output")
}

func (r *Root) initializeFormatFlags() {
	r.Command.Flags().IntVarP(
		&r.Options.Indent,
		"indent-level",
		"I",
		2,
		helpIndentLevel,
	)

	r.Command.Flags().VarP(
		enumflag.NewSlice(&r.Options.Styles, "style", actions.Styles{}.FlagMap(), enumflag.EnumCaseInsensitive),
		"output-style",
		"S",
		helpOutputStyle,
	)

	r.Command.Flags().BoolVarP(
		&r.Options.RemoveComments,
		"remove-comments",
		"",
		false,
		helpRemoveComments,
	)

	if err := viper.BindPFlag("removeComments", r.Command.Flags().Lookup("remove-comments")); err != nil {
		r.Log.Fatal(err)
	}

	if err := viper.BindPFlag("indentLevel", r.Command.Flags().Lookup("indent-level")); err != nil {
		r.Log.Fatal(err)
	}

	if err := viper.BindPFlag("outputStyle", r.Command.Flags().Lookup("output-style")); err != nil {
		r.Log.Fatal(err)
	}
}

func (r *Root) initializeTemplateFlags() {
	r.Command.Flags().StringArrayVarP(
		&r.Options.ValueFiles,
		"values-file",
		"f",
		nil,
		helpValueFile,
	)
}

func (r *Root) initializeInstructionFlags() {
	r.Command.Flags().StringVarP(
		&r.Options.InstructionsFile,
		"instructions",
		"i",
		"",
		helpInstructionsFile,
	)

	if err := r.Command.MarkFlagFilename("instructions"); err != nil {
		r.Log.Error(err)
	}

	r.Command.Flags().Var(
		&r.Options.DefaultOnMissingAction,
		"default-on-missing-action",
		helpDefaultOnMissingAction,
	)

	if err := viper.BindPFlag("defaultOnMissingAction", r.Command.Flags().Lookup("default-on-missing-action")); err != nil {
		r.Log.Fatal(err)
	}
}

func (r *Root) initializeStdInFlags() {
	r.Command.Flags().VarP(
		&r.Options.Overlay.Query,
		"query",
		"q",
		helpQuery,
	)

	r.Command.Flags().VarP(
		&r.Options.Overlay.Action,
		"action",
		"a",
		helpAction,
	)

	r.Command.Flags().StringVarP(
		&r.Options.Value,
		"value",
		"x",
		"",
		helpValue,
	)
	r.Command.Flags().StringVarP(
		&r.Options.Path,
		"path",
		"p",
		"",
		helpPath,
	)
}
