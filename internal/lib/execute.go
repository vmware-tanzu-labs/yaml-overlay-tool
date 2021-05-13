// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: MIT

package lib

import (
	"errors"

	"github.com/op/go-logging"
)

var (
	log              = logging.MustGetLogger("lib") //nolint:gochecknoglobals
	ErrInvalidAction = errors.New("invalid overlay action")
)

func Execute(options *Options) error {
	instructions, err := ReadInstructionFile(&options.InstructionsFile)
	if err != nil {
		return err
	}

	return instructions.processYamlFiles(options)
}
