package commands

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/vmware-tanzu-labs/yaml-overlay-tool/internal/lib"
)

var options lib.Options //nolint:gochecknoglobals

func initializeGlobalFlags(rootCmd *cobra.Command) {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().BoolVarP(&options.Verbose, "verbose", "V", false, "verbose mode")
	rootCmd.PersistentFlags().StringP(
		"config",
		"c",
		"~/.yot",
		"config file (default is $HOME/.yot)",
	)

	// Define with Cobra
	rootCmd.PersistentFlags().StringP(
		"common-values",
		"d",
		"",
		helpCommonValues,
	)

	rootCmd.PersistentFlags().StringP(
		"default-values-file",
		"",
		"",
		helpDefaultValueFile,
	)

	err := rootCmd.PersistentFlags().MarkDeprecated("default-values-file", helpDefaultValuesFileDeprecated)
	if err != nil {
		log.Fatal(err)
	}

	rootCmd.PersistentFlags().StringSliceVarP(
		&options.ValuesPath,
		"values",
		"f",
		[]string{},
		helpValuesPath,
	)

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

	rootCmd.PersistentFlags().StringP(
		"dump-rendered-instructions",
		"r",
		"",
		helpDumpRenderedInstructions,
	)

	rootCmd.PersistentFlags().IntVarP(
		&options.Indent,
		"indent-level",
		"I",
		2,
		helpIndentLevel,
	)

	// Bind w/ viper
	viper.BindPFlag("default-values-file", rootCmd.PersistentFlags().Lookup("default-values-file"))
	viper.BindPFlag("values-path", rootCmd.PersistentFlags().Lookup("values-path"))
	viper.BindPFlag("instruction-file", rootCmd.PersistentFlags().Lookup("instruction-file"))
	viper.BindPFlag("output-directory", rootCmd.PersistentFlags().Lookup("output-directory"))
	viper.BindPFlag("stdout", rootCmd.PersistentFlags().Lookup("stdout"))
	viper.BindPFlag("dump-rendered-instructions", rootCmd.PersistentFlags().Lookup("dump-rendered-instructions"))
}
