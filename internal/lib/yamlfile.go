// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: MIT

package lib

import (
	"fmt"
)

type Document struct {
	Name     string     `yaml:"name,omitempty"`
	Path     int        `yaml:"index,omitempty"`
	Overlays []*Overlay `yaml:"overlays,omitempty"`
}

type YamlFile struct {
	Name      string      `yaml:"name,omitempty"`
	Overlays  []*Overlay  `yaml:"overlays,omitempty"`
	Documents []*Document `yaml:"documents,omitempty"`
	Sources   Sources     `yaml:"path,omitempty"`
}

func (yf *YamlFile) processYamlFiles(cfg *Config) error {
	for _, src := range yf.Sources {
		for nodeIndex := range src.Nodes {
			log.Infof("Processing Common & File Overlays in File %s on Document %d\n\n", src.Path, nodeIndex)

			if err := src.processOverlays(yf.Overlays, nodeIndex); err != nil {
				return fmt.Errorf("failed to apply file overlays, %w", err)
			}

			log.Infof("Processing Document Overlays in File %s on Document %d\n\n", src.Path, nodeIndex)

			for di := range yf.Documents {
				if ok := yf.Documents[di].checkDocumentPath(nodeIndex); !ok {
					continue
				}

				if err := src.processOverlays(yf.Documents[di].Overlays, nodeIndex); err != nil {
					return err
				}
			}
		}

		if err := src.doPostProcessing(cfg); err != nil {
			return fmt.Errorf("failed to perform post processing on %s: %w", src.Path, err)
		}
	}

	return nil
}

func (d *Document) checkDocumentPath(docIndex int) bool {
	return d.Path == docIndex
}
