// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: MIT

package lib

type Document struct {
	Name     string     `yaml:"name,omitempty"`
	Path     int        `yaml:"index,omitempty"`
	Overlays []*Overlay `yaml:"overlays,omitempty"`
}

func (d *Document) checkDocumentIndex(docIndex int) bool {
	return d.Path == docIndex
}
