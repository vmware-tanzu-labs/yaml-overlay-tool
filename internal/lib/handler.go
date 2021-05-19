// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: MIT

package lib

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

type workStream struct {
	Overlay   Overlay
	Node      *yaml.Node
	NodeIndex int
	File      *File
}

func OverlayHandler(cfg *Config, oChan chan *workStream, errs chan error) {
	for o := range oChan {
		log.Noticef("Processing overlay [%q] on document %d of file %s", o.Overlay.Name, o.NodeIndex, o.File.Path)

		if err := o.Overlay.apply(o.Node); err != nil {
			errs <- fmt.Errorf("%w in file %s on Document %d", err, o.File.Path, o.NodeIndex)
		}
	}
}
