// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: MIT

package lib

import (
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("lib") //nolint:gochecknoglobals

func Execute(options *Options) error {
	instructions, err := ReadInstructionFile(&options.InstructionsFile)
	if err != nil {
		return err
	}

	return instructions.processYamlFiles(options)
}
