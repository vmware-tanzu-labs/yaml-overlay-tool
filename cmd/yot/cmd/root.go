/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"errors"
	"fmt"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "yot [-h] [-d DEFAULT_VALUES_FILE] [-v VALUES_PATH] -i INSTRUCTION_FILE [-o OUTPUT_DIRECTORY] [-s] [-r] [-l LOG_FILE] [-V]",
	Short: "yot (YAML overlay tool) is a yaml overlay tool which allows for the templating of overlay instruction data with jinja2",
	Long: `yot (YAML overlay tool) is a yaml overlay tool which allows for the templating 
of overlay instruction data with jinja2, and the application of rendered 
overlays "over the top" of a yaml file. yot only produces valid yaml 
documents on output.`,
	Version: "yaml overlay tool v0.01",
	// PreRun check for default usage requirements and run example
	// PreRunE: func(cmd *cobra.Command, args []string) error {
	// 	return CheckRequiredFlags(cmd.Flags())
	// },
	// Uncomment the following line if your bare application
	// has an action associated with it:
	RunE: func(cmd *cobra.Command, args []string) error {
		// TODO: this needs to be fixed so that it passed by ref and NOT to use the OS obejct.
		configPath := os.Args[2]
		if err := instructFile(configPath); err != nil {
			return err
		}
		os.Exit(0)
		return nil
	},
	Args: func(cmd *cobra.Command, args []string) error {
		// TODO: this needs to be fixed so that it passed by ref and NOT to use the OS object.
		fmt.Println(len(os.Args), os.Args)
		fmt.Println(len(args), args)
		if len(os.Args) < 1 {
			return errors.New("requires input single file (.yaml/.yml) or directory with yaml files ")
		}
		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.yot.yaml)")

	// Define with Cobra
	rootCmd.PersistentFlags().StringP("default-values-file", "d", "", "Path to your default values file. If not set, you must \npass a values file of defaults.yaml or \ndefaults.yml within a path from the -v option. \nTakes multiple default values files in case you would \nlike to separate out some of the values. After the \nfirst default values file, each subsequent file passed \nwith -d will be merged with the values from the \nfirst. If a defaults.yaml or defaults.yml file is \ndiscovered in one of your -v paths, it will be \nmerged with these values last.")
	rootCmd.PersistentFlags().StringP("values-path", "v", "", "Values file path. May be a path to a file or directory \ncontaining value files ending in either .yml or .yaml. \nThis option can be provided multiple times \nas required. A file named defaults.yaml or \ndefaults.yml is required within the path(s) if not \nusing the -d option, and you may have only 1 default \nvalue file in that scenario. Additional values files \nare merged over the defaults.yaml file values. Each \nvalues file is treated as a unique site and will \nrender your instructions differently based on its \nvalues")
	rootCmd.PersistentFlags().StringP("instruction-file", "i", "", "Instruction file path. Defaults to ./instructions.yaml (required)")
	rootCmd.MarkPersistentFlagRequired("instruction-file")
	rootCmd.PersistentFlags().StringP("output-directory", "o", "", "Path to directory to write the overlayed yaml files \nto. If value files were supplied in addition to a \ndefaults.yaml/.yml then the rendered templates will \nland in <output dir>/<addl value file name>.")
	rootCmd.PersistentFlags().StringP("stdout", "s", "", "Render output to stdout. Templated instructions files \nwill still be output to the --output-directory.")
	rootCmd.PersistentFlags().StringP("dump-rendered-instructions", "r", "", "If using a templated instructions file, you can dump \nthe rendered instructions to stdout to allow for \nreviewing how they were rendered prior to a full run \nof yot. Equivalent to a dry-run. Exits with return \ncode 0 prior to processing instructions")

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

// initConfig reads in config file and ENV variables if set.
func initConfig() {
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
