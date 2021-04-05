package commands

import (
	"github.com/spf13/cobra"
)

// versionCmd represents the version command
// TODO: This cammand can be used to list all availible version of this tool and potentilly invoke and
// download newer version.  The essance of this sub command can be functional especially if the version
// is associated with a COMMIT hash in github.
var versionCmd = &cobra.Command{
	Use:     "version",
	Short:   "yaml overlay tool v0.0.1",
	Long:    `yaml overlay tool v0.0.1`,
	Version: "yaml overlay tool v0.0.1",
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
