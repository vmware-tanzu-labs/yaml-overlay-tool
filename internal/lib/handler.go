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
	Path      string
}

func OverlayHandler(cfg *Config, oChan chan *workStream, errs chan error) {
	for o := range oChan {
		log.Noticef("Processing overlay [%q] in file %s on document %d", o.Overlay.Name, o.Path, o.NodeIndex)

		if err := o.Overlay.apply(o.Node); err != nil {
			errs <- fmt.Errorf("%w in file %s on Document %d", err, o.Path, o.NodeIndex)
		}
	}
}
