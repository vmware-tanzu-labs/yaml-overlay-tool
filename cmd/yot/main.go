package main

import (
	"os"

	command "github.com/vmware-tanzu-labs/yaml-overlay-tool/internal/commands"
)

func main() {
	cmd := command.New()
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
