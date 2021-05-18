// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: MIT

package lib

type Config struct {
	Verbose          bool
	LogLevel         string
	InstructionsFile string
	OutputDir        string
	StdOut           bool
	Indent           int
}
