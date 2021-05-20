// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: MIT

package instructions

type Config struct {
	Verbose          bool
	LogLevel         string
	InstructionsFile string
	OutputDir        string
	StdOut           bool
	Indent           int
}
