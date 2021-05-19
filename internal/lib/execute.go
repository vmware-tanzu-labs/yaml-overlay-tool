// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: MIT

package lib

import (
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("lib") //nolint:gochecknoglobals

func Execute(cfg *Config) error {
	errs := make(chan error)

	instructions, err := ReadInstructionFile(&cfg.InstructionsFile)
	if err != nil {
		return err
	}

	instructions.addCommonOverlays()

	yfChan := make(chan *YamlFile, len(instructions.YamlFiles))

	go instructions.queueYamlFiles(yfChan)

	go func() {
		for yf := range yfChan {
			oChan := make(chan *workStream)

			go yf.queueOverlays(oChan)

			OverlayHandler(cfg, oChan, errs)

			if err := yf.doPostProcessing(cfg); err != nil {
				errs <- err
			}
		}

		close(errs)
	}()

	if err, ok := <-errs; ok && err != nil {
		return err
	}

	return nil
}
