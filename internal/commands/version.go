package commands

import (
	"github.com/spf13/cobra"
)

func createVersionCommand() *cobra.Command {
	// versionCmd represents the version command
	// TODO: This cammand can be used to list all availible version of this tool and potentilly invoke and
	// download newer version.  The essance of this sub command can be functional especially if the version
	// is associated with a COMMIT hash in github.
	return &cobra.Command{
		Use:     "version",
		Short:   "yaml overlay tool v0.0.1",
		Long:    `yaml overlay tool v0.0.1`,
		Version: "yaml overlay tool v0.0.1",
		// Run: func(cmd *cobra.Command, args []string) {
		// fmt.Println("yaml overlay tool version v0.01")
		// },
	}
}
