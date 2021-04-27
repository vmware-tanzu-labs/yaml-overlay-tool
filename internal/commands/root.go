// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: MIT

package commands

import (
	"fmt"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/op/go-logging"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/vmware-tanzu-labs/yaml-overlay-tool/internal/lib"
)

var ErrMissingRequired = fmt.Errorf("missing required arguments")

func New() *cobra.Command {
	// rootCmd represents the base command when called without any subcommands.
	rootCmd := &cobra.Command{
		Use:     yotUsage,
		Short:   yotShort,
		Long:    yotLong,
		Version: "yaml overlay tool v0.0.1",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := lib.Execute(&options); err != nil {
				cmd.SilenceUsage = true

				return fmt.Errorf("%w", err)
			}

			return nil
		},
		PersistentPreRun: SetupLogging,
	}

	cobra.OnInitialize(initConfig)

	initializeGlobalFlags(rootCmd)

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
		return fmt.Errorf("%w: %q has not been set", ErrMissingRequired, flagName)
	}

	return nil
}

func SetupLogging(cmd *cobra.Command, args []string) {
	cmd.SetOut(cmd.OutOrStdout())

	format := logging.MustStringFormatter(
		`%{color}%{time:15:04:05} [%{level}]%{color:reset} %{message}`,
	)
	backend := logging.AddModuleLevel(
		logging.NewBackendFormatter(logging.NewLogBackend(os.Stderr, "", 0), format))

	if options.Verbose {
		backend.SetLevel(logging.DEBUG, "")
	} else {
		backend.SetLevel(logging.ERROR, "")
	}

	logging.SetBackend(backend)
}
