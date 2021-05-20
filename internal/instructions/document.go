// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: MIT

package instructions

import "github.com/vmware-tanzu-labs/yaml-overlay-tool/internal/overlays"

type Document struct {
	Name     string              `yaml:"name,omitempty"`
	Path     int                 `yaml:"index,omitempty"`
	Overlays []*overlays.Overlay `yaml:"overlays,omitempty"`
}

func (d *Document) checkDocumentIndex(docIndex int) bool {
	return d.Path == docIndex
}
