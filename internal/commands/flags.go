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
		"default-values-file",
		"d",
		"",
		helpDefaultValueFile,
	)

	rootCmd.PersistentFlags().StringSliceP(
		"values-path",
		"v",
		[]string{},
		helpValuesPath,
	)

	rootCmd.PersistentFlags().StringVarP(
		&options.InstructionsFile,
		"instruction-file",
		"i",
		"",
		helpInstructionsFile,
	)

	if err := rootCmd.MarkPersistentFlagRequired("instruction-file"); err != nil {
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

	// Bind w/ viper
	viper.BindPFlag("default-values-file", rootCmd.PersistentFlags().Lookup("default-values-file"))
	viper.BindPFlag("values-path", rootCmd.PersistentFlags().Lookup("values-path"))
	viper.BindPFlag("instruction-file", rootCmd.PersistentFlags().Lookup("instruction-file"))
	viper.BindPFlag("output-directory", rootCmd.PersistentFlags().Lookup("output-directory"))
	viper.BindPFlag("stdout", rootCmd.PersistentFlags().Lookup("stdout"))
	viper.BindPFlag("dump-rendered-instructions", rootCmd.PersistentFlags().Lookup("dump-rendered-instructions"))

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
