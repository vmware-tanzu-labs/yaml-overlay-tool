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
	r.Command.Flags().StringVar(&r.configFile, "config", "", "config file (default is $HOME/.yot)")

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

	if err := viper.BindPFlag("log-level", r.Command.Flags().Lookup("log-level")); err != nil {
		r.Log.Fatal(err)
	}

	viper.RegisterAlias("logLevel", "log-level")
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

	if err := viper.BindPFlag("output-directory", r.Command.Flags().Lookup("output-directory")); err != nil {
		r.Log.Fatal(err)
	}

	viper.RegisterAlias("outputDirectory", "output-directory")

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

	if err := viper.BindPFlag("remove-comments", r.Command.Flags().Lookup("remove-comments")); err != nil {
		r.Log.Fatal(err)
	}

	if err := viper.BindPFlag("indent-level", r.Command.Flags().Lookup("indent-level")); err != nil {
		r.Log.Fatal(err)
	}

	if err := viper.BindPFlag("output-style", r.Command.Flags().Lookup("output-style")); err != nil {
		r.Log.Fatal(err)
	}

	viper.RegisterAlias("indentLevel", "indent-level")
	viper.RegisterAlias("outputStyle", "output-style")
	viper.RegisterAlias("removeComments", "remove-comments")
}

func (r *Root) initializeTemplateFlags() {
	r.Command.Flags().StringArrayVarP(
		&r.Options.Values,
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

	if err := viper.BindPFlag("default-on-missing-action", r.Command.Flags().Lookup("default-on-missing-action")); err != nil {
		r.Log.Fatal(err)
	}

	viper.RegisterAlias("defaultOnMissingAction", "default-on-missing-action")
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
