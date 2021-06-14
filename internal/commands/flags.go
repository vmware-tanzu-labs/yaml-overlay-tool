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
		helpLogLevel,
	)

	r.Command.Flags().Lookup("log-level").NoOptDefVal = "debug"

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

	r.Command.Flags().BoolVarP(
		&r.Options.RemoveComments,
		"remove-comments",
		"",
		false,
		helpRemoveComments,
	)

	r.Command.Flags().IntVarP(
		&r.Options.Indent,
		"indent-level",
		"I",
		2,
		helpIndentLevel,
	)

	r.Command.Flags().VarP(
		enumflag.NewSlice(&r.Options.Styles, "style", r.Options.Styles.FlagMap(), enumflag.EnumCaseInsensitive),
		"output-style",
		"S",
		helpOutputStyle,
	)
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
