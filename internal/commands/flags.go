// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: MIT

package commands

import (
	"github.com/op/go-logging"
	"github.com/thediveo/enumflag"
)

func (r *Root) initializeGlobalFlags() {
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
		HelpLogLevel,
	)

	r.Command.Flags().Lookup("log-level").NoOptDefVal = "debug"

	r.Command.Flags().StringVarP(
		&r.Options.InstructionsFile,
		"instructions",
		"i",
		"instructions.yaml",
		HelpInstructionsFile,
	)

	if err := r.Command.MarkFlagFilename("instructions"); err != nil {
		r.Log.Error(err)
	}

	r.Command.Flags().StringVarP(
		&r.Options.OutputDir,
		"output-directory",
		"o",
		"./output",
		HelpOutputDirectory,
	)

	if err := r.Command.MarkFlagDirname("output-directory"); err != nil {
		r.Log.Fatal(err)
	}

	r.Command.Flags().BoolVarP(
		&r.Options.StdOut,
		"stdout",
		"s",
		false,
		HelpRenderStdOut,
	)

	r.Command.Flags().BoolVarP(
		&r.Options.RemoveComments,
		"remove-comments",
		"",
		false,
		HelpRemoveComments,
	)

	r.Command.Flags().IntVarP(
		&r.Options.Indent,
		"indent-level",
		"I",
		2,
		HelpIndentLevel,
	)

	r.Command.Flags().VarP(
		enumflag.NewSlice(&r.Options.Styles, "style", r.Options.Styles.FlagMap(), enumflag.EnumCaseInsensitive),
		"output-style",
		"S",
		HelpOutputStyle,
	)
}

func (r *Root) initializeTemplateFlags() {
	r.Command.Flags().StringArrayVarP(
		&r.Options.Values,
		"values",
		"f",
		nil,
		helpValueFile,
	)
}
