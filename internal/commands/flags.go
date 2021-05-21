// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: MIT

package commands

import (
	"github.com/op/go-logging"
	"github.com/spf13/cobra"
	"github.com/thediveo/enumflag"
)

func (rc *RootCommand) initializeGlobalFlags(rootCmd *cobra.Command) {
	logMap := map[logging.Level][]string{
		logging.CRITICAL: {"critical", "crit", "c"},
		logging.ERROR:    {"error", "err", "e"},
		logging.WARNING:  {"warning", "warn", "w"},
		logging.NOTICE:   {"notice", "note", "n"},
		logging.INFO:     {"info", "i"},
		logging.DEBUG:    {"debug", "d", "verbose", "v"},
	}

	rootCmd.Flags().VarP(
		enumflag.New(&rc.Options.LogLevel, "logLevel", logMap, enumflag.EnumCaseInsensitive),
		"log-level",
		"v",
		HelpLogLevel,
	)

	rootCmd.Flags().Lookup("log-level").NoOptDefVal = "debug"

	rootCmd.Flags().StringVarP(
		&rc.Options.InstructionsFile,
		"instructions",
		"i",
		"instructions.yaml",
		HelpInstructionsFile,
	)

	if err := rootCmd.MarkFlagFilename("instructions"); err != nil {
		rc.Log.Error(err)
	}

	rootCmd.Flags().StringVarP(
		&rc.Options.OutputDir,
		"output-directory",
		"o",
		"./output",
		HelpOutputDirectory,
	)

	if err := rootCmd.MarkFlagDirname("output-directory"); err != nil {
		rc.Log.Fatal(err)
	}

	rootCmd.Flags().BoolVarP(
		&rc.Options.StdOut,
		"stdout",
		"s",
		false,
		HelpRenderStdOut,
	)

	rootCmd.Flags().IntVarP(
		&rc.Options.Indent,
		"indent-level",
		"I",
		2,
		HelpIndentLevel,
	)

	rootCmd.Flags().VarP(
		enumflag.NewSlice(&rc.Options.Styles, "style", rc.Options.Styles.FlagMap(), enumflag.EnumCaseInsensitive),
		"output-style",
		"S",
		HelpOutputStyle,
	)

	rootCmd.Flags()
}
