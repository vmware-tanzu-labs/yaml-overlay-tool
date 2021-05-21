// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: MIT

package instructions

import (
	"github.com/op/go-logging"
	"github.com/vmware-tanzu-labs/yaml-overlay-tool/internal/actions"
)

// Config contains configuration options used with instruction files.
type Config struct {
	Verbose          bool
	LogLevel         logging.Level
	InstructionsFile string
	OutputDir        string
	StdOut           bool
	Indent           int
	Styles           actions.Styles
}
