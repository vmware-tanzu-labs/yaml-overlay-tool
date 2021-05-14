// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: MIT

package lib

import (
	"errors"
	"fmt"
	"io"
	"os"

	"gopkg.in/yaml.v3"
)

type Document struct {
	Name     string    `yaml:"name,omitempty"`
	Path     string    `yaml:"path,omitempty"`
	Overlays []Overlay `yaml:"overlays,omitempty"`
}

type YamlFile struct {
	Document  `yaml:",inline"`
	Documents []Document `yaml:"documents,omitempty"`
	Source    []Source
}

func (yf *YamlFile) readYamlFile(path string) error {
	source := &Source{
		Origin: "file",
		Path:   path,
	}

	reader, err := ReadStream(path)
	if err != nil {
		return err
	}

	dc := yaml.NewDecoder(reader)

	for {
		var y yaml.Node

		if err := dc.Decode(&y); errors.Is(err, io.EOF) {
			if reader, ok := reader.(*os.File); ok {
				CloseFile(reader)

				break
			}
		} else if err != nil {
			return fmt.Errorf("failed to read file %s: %w", yf.Path, err)
		}

		source.Nodes = append(source.Nodes, &y)
	}

	yf.Source = append(yf.Source, *source)

	return nil
}

func (yf *YamlFile) processYamlFiles(opt *Options) error {
	for _, src := range yf.Source {
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

		if err := src.doPostProcessing(opt); err != nil {
			return fmt.Errorf("failed to perform post processing on %s: %w", src.Path, err)
		}
	}

	return nil
}

func (d *Document) checkDocumentPath(docIndex int) bool {
	return d.Path == fmt.Sprint(docIndex)
}
