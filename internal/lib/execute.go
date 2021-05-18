// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: MIT

package lib

import (
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("lib") //nolint:gochecknoglobals

func Execute(cfg *Config) error {
	instructions, err := ReadInstructionFile(&cfg.InstructionsFile)
	if err != nil {
		return err
	}

	instructions.addCommonOverlays()

	for yfIndex := range instructions.YamlFiles {
		if err := instructions.YamlFiles[yfIndex].processYamlFiles(cfg); err != nil {
			return err
		}
	}

	return nil
}
