// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: MIT

package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func (r *Root) EnvCommand() *cobra.Command {
	return &cobra.Command{
		Use:                   "env [settingName]",
		Short:                 envShort,
		Long:                  envLong,
		DisableFlagsInUseLine: true,
		Args:                  cobra.MaximumNArgs(1),
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			if len(args) == 0 {
				kmap := envMap()
				keys := make([]string, len(kmap))
				i := 0
				for k := range kmap {
					keys[i] = k
					i++
				}

				return keys, cobra.ShellCompDirectiveNoFileComp
			}

			return nil, cobra.ShellCompDirectiveNoFileComp
		},
		Run: r.env,
	}
}

func (r *Root) env(cmd *cobra.Command, args []string) {
	settings := envMap()
	if len(args) == 0 {
		for k, v := range settings {
			os.Stdout.WriteString(fmt.Sprintf("%s=\"%v\"\n", k, viper.Get(v)))
		}
	} else {
		os.Stdout.WriteString(fmt.Sprintf("%s=\"%v\"\n", args[0], viper.Get(settings[args[0]])))
	}
}

func envMap() map[string]string {
	return map[string]string{
		"YOT_CONFIG_FILE":               "configFile",
		"YOT_DEFAULT_ON_MISSING_ACTION": "defaultOnMissingAction",
		"YOT_INDENT_LEVEL":              "indentLevel",
		"YOT_LOG_LEVEL":                 "logLevel",
		"YOT_OUTPUT_DIRECTORY":          "outputDirectory",
		"YOT_OUTPUT_STYLE":              "outputStyle",
		"YOT_REMOVE_COMMENTS":           "removeComments",
		"YOT_STDOUT":                    "stdout",
	}
}
