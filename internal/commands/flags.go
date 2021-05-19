// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: MIT

package commands

import (
	"github.com/spf13/cobra"
)

func (rc *RootCommand) initializeGlobalFlags(rootCmd *cobra.Command) {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.Flags().BoolVarP(
		&rc.Options.Verbose,
		"verbose",
		"V",
		false,
		HelpVerbose,
	)

	rootCmd.Flags().StringVarP(
		&rc.Options.LogLevel,
		"log-level",
		"l",
		"",
		HelpLogLevel,
	)

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
}
