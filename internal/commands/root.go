// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: MIT

package commands

import (
	"errors"
	"fmt"
	"log"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/op/go-logging"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/vmware-tanzu-labs/yaml-overlay-tool/internal/lib"
)

func New() *cobra.Command {

	// rootCmd represents the base command when called without any subcommands.
	var rootCmd = &cobra.Command{
		Use:   "yot [-h] [-d DEFAULT_VALUES_FILE] [-v VALUES_PATH] -i INSTRUCTION_FILE [-o OUTPUT_DIRECTORY] [-s] [-r] [-l LOG_FILE] [-V]",
		Short: "yot (YAML overlay tool) is a yaml overlay tool which allows for the templating of overlay instruction data with jinja2",
		Long: `yot (YAML overlay tool) is a yaml overlay tool which allows for the templating 
	of overlay instruction data with jinja2, and the application of rendered 
	overlays "over the top" of a yaml file. yot only produces valid yaml 
	documents on output.`,
		Version: "yaml overlay tool v0.0.1",
		RunE: func(cmd *cobra.Command, args []string) error {
			instructions, err := lib.ReadInstructionFile(&instructionFile)
			if err != nil {
				return err
			}

			return lib.Process(instructions)
		},
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			cmd.SetOut(cmd.OutOrStdout())
			var format = logging.MustStringFormatter(
				`%{color}%{time:15:04:05} [%{level}]%{color:reset} %{message}`,
			)
			var backend = logging.AddModuleLevel(
				logging.NewBackendFormatter(logging.NewLogBackend(os.Stderr, "", 0), format))

			if verbose {
				backend.SetLevel(logging.DEBUG, "")
			} else {
				backend.SetLevel(logging.ERROR, "")
			}

			logging.SetBackend(backend)
		},
	}

	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "V", false, "verbose mode")
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
		`Path to your default values file. If not set, you must 
pass a values file of defaults.yaml or 
defaults.yml within a path from the -v option. 
Takes multiple default values files in case you would 
like to separate out some of the values. After the 
first default values file, each subsequent file passed 
with -d will be merged with the values from the 
first. If a defaults.yaml or defaults.yml file is 
discovered in one of your -v paths, it will be 
merged with these values last.`,
	)

	rootCmd.PersistentFlags().StringSliceP(
		"values-path",
		"v",
		[]string{},
		`Values file path. May be a path to a file or directory 
containing value files ending in either .yml or .yaml. 
This option can be provided multiple times as required. 
A file named defaults.yaml or defaults.yml is required 
within the path(s) if not using the -d option, and you
may have only 1 default value file in that scenario. 
Additional values files are merged over the defaults.yaml
file values. Each values file is treated as a unique site
and will render your instructions differently based on its
values`,
	)

	rootCmd.PersistentFlags().StringVarP(
		&instructionFile,
		"instruction-file",
		"i",
		"",
		"Instruction file path. Defaults to ./instructions.yaml (required)",
	)

	if err := rootCmd.MarkPersistentFlagRequired("instruction-file"); err != nil {
		log.Fatal("InstructionsFile (-i) is required")
	}

	rootCmd.PersistentFlags().StringP(
		"output-directory",
		"o",
		"./output",
		`Path to directory to write the overlayed yaml files to.
If value files were supplied in addition to a 
defaults.yaml/.yml then the rendered templates will land
in <output dir>/<addl value file name>.`,
	)

	rootCmd.PersistentFlags().StringP(
		"stdout",
		"s",
		"",
		`Render output to stdout. Templated instructions files 
will still be output to the --output-directory.`,
	)

	rootCmd.PersistentFlags().StringP(
		"dump-rendered-instructions",
		"r",
		"",
		`If using a templated instructions file, you can dump 
the rendered instructions to stdout to allow for 
reviewing how they were rendered prior to a full run 
of yot. Equivalent to a dry-run. Exits with return
code 0 prior to processing instructions`,
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

	rootCmd.AddCommand(createVersionCommand())

	return rootCmd
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	var cfgFile string

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".yot" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".yot")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}

// Check for required flags identified by MarkPersistentFlagRequired.
func CheckRequiredFlags(flags *pflag.FlagSet) error {
	requiredError := false
	flagName := ""

	flags.VisitAll(func(flag *pflag.Flag) {
		requiredAnnotation := flag.Annotations[cobra.BashCompOneRequiredFlag]
		if len(requiredAnnotation) == 0 {
			return
		}

		flagRequired := requiredAnnotation[0] == "true"

		if flagRequired && !flag.Changed {
			requiredError = true
			flagName = flag.Name
		}
	})

	if requiredError {
		return errors.New("yot: error: the following arguments are required: `" + flagName + "` has not been set")
	}

	return nil
}
