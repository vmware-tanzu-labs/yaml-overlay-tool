// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: MIT

package commands

import (
	"github.com/spf13/cobra"
	"github.com/vmware-tanzu-labs/yaml-overlay-tool/internal/lib"
)

var options lib.Options //nolint:gochecknoglobals

func initializeGlobalFlags(rootCmd *cobra.Command) {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().BoolVarP(
		&options.Verbose,
		"verbose",
		"V",
		false,
		HelpVerbose,
	)

	rootCmd.PersistentFlags().StringVarP(
		&options.LogLevel,
		"log-level",
		"l",
		"",
		HelpLogLevel,
	)

	rootCmd.Flags().StringVarP(
		&options.InstructionsFile,
		"instructions",
		"i",
		"instructions.yaml",
		HelpInstructionsFile,
	)

	if err := rootCmd.MarkFlagFilename("instructions"); err != nil {
		log.Error(err)
	}

	rootCmd.PersistentFlags().StringVarP(
		&options.OutputDir,
		"output-directory",
		"o",
		"./output",
		HelpOutputDirectory,
	)

	if err := rootCmd.MarkPersistentFlagDirname("output-directory"); err != nil {
		log.Fatal(err)
	}

	rootCmd.PersistentFlags().BoolVarP(
		&options.StdOut,
		"stdout",
		"s",
		false,
		HelpRenderStdOut,
	)

	rootCmd.PersistentFlags().IntVarP(
		&options.Indent,
		"indent-level",
		"I",
		2,
		HelpIndentLevel,
	)
}
