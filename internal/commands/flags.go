// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: MIT

package commands

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/vmware-tanzu-labs/yaml-overlay-tool/internal/lib"
)

var options lib.Options //nolint:gochecknoglobals

func initializeGlobalFlags(rootCmd *cobra.Command) {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().BoolVarP(&options.Verbose, "verbose", "V", false, "verbose mode")

	rootCmd.PersistentFlags().StringVarP(
		&options.InstructionsFile,
		"instructions",
		"i",
		"",
		helpInstructionsFile,
	)

	if err := rootCmd.MarkPersistentFlagRequired("instructions"); err != nil {
		log.Fatal("InstructionsFile (-i) is required")
	}

	rootCmd.PersistentFlags().StringVarP(
		&options.OutputDir,
		"output-directory",
		"o",
		"./output",
		helpOutputDirectory,
	)

	rootCmd.PersistentFlags().BoolVarP(
		&options.StdOut,
		"stdout",
		"s",
		false,
		helpRenderStdOut,
	)

	rootCmd.PersistentFlags().IntVarP(
		&options.Indent,
		"indent-level",
		"I",
		2,
		helpIndentLevel,
	)
}
