// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: MIT

package instructions

import "github.com/vmware-tanzu-labs/yaml-overlay-tool/internal/overlays"

// Document is a set of overlays scoped to a specific document within a yamlfile.
type Document struct {
	Name     string              `yaml:"name,omitempty"`
	Path     int                 `yaml:"path,omitempty"`
	Overlays []*overlays.Overlay `yaml:"overlays,omitempty"`
}

// checkDocumentIndex checks to see if the current index matches a Document's path.
func (d *Document) checkDocumentIndex(docIndex int) bool {
	return d.Path == docIndex
}
