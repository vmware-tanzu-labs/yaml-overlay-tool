// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: MIT

package main

import (
	command "github.com/vmware-tanzu-labs/yaml-overlay-tool/internal/commands"
)

var version = "unstable"

func main() {
	yot := command.New().Command(version)

	yot.Execute() //nolint:errcheck // not needed as command will handle errors
}
