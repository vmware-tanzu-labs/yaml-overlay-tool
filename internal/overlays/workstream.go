// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: MIT

package overlays

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

type workload struct {
	Overlay   *Overlay
	Node      *yaml.Node
	NodeIndex int
	Path      string
}

type WorkStream struct {
	stream chan *workload
}

func NewWorkStream() *WorkStream {
	return &WorkStream{}
}

func (s *WorkStream) StartStream() {
	s.stream = make(chan *workload)
}

func (s *WorkStream) AddWorkload(o *Overlay, n *yaml.Node, nodeIndex int, path string) {
	s.stream <- &workload{
		Overlay:   o,
		Node:      n,
		NodeIndex: nodeIndex,
		Path:      path,
	}
}

func (s *WorkStream) CloseStream() {
	close(s.stream)
}

func (s *WorkStream) StartHandler() error {
	for o := range s.stream {
		log.Noticef("Processing overlay [%q] in file %s on document %d", o.Overlay.Name, o.Path, o.NodeIndex)

		if err := o.Overlay.Apply(o.Node); err != nil {
			return fmt.Errorf("%w in file %s on Document %d", err, o.Path, o.NodeIndex)
		}
	}

	return nil
}
