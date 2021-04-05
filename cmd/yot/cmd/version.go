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
	"github.com/spf13/cobra"
)

// versionCmd represents the version command
// TODO: This cammand can be used to list all availible version of this tool and potentilly invoke and
// download newer version.  The essance of this sub command can be functional especially if the version
// is associated with a COMMIT hash in github
var versionCmd = &cobra.Command{
	Use:     "version",
	Short:   "yaml overlay tool v0.01",
	Long:    `yaml overlay tool v0.01`,
	Version: "yaml overlay tool v0.01",
	// Run: func(cmd *cobra.Command, args []string) {
	// fmt.Println("yaml overlay tool version v0.01")
	// },
}

func init() {
	rootCmd.AddCommand(versionCmd)
	// Here you will define your flags and configuration settings.
	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// versionCmd.PersistentFlags().String("foo", "", "A help for foo")
	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// versionCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
