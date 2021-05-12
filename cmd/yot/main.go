// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: MIT

package main

import (
	"log"

	command "github.com/vmware-tanzu-labs/yaml-overlay-tool/internal/commands"
)

func main() {
	yot := command.New()

	if err := yot.Execute(); err != nil {
		log.Fatal(err)
	}
}
