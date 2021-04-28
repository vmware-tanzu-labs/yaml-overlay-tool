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

func ReadInstructionFile(fileName *string) (*Instructions, error) {
	var instructions Instructions

	log.Debugf("Instructions File: %s\n\n", *fileName)

	reader, err := ReadStream(*fileName)
	if err != nil {
		return nil, err
	}

	dc := yaml.NewDecoder(reader)
	if err := dc.Decode(&instructions); err != nil {
		return nil, fmt.Errorf("unable to read instructions file %s: %w", *fileName, err)
	}

	if err := instructions.readYamlFiles(); err != nil {
		return nil, err
	}

	return &instructions, nil
}

func (i *Instructions) applyOverlays(options *Options) error {
	for _, file := range i.YamlFiles {
		for nodeIndex := range file.Nodes {
			log.Infof("Processing Common Overlays in File %s on Document %d\n\n", file.Path, nodeIndex)

			if err := file.processOverlays(i.CommonOverlays, nodeIndex); err != nil {
				return fmt.Errorf("failed to apply common overlays, %w", err)
			}

			log.Infof("Processing File Overlays in File %s on Document %d\n\n", file.Path, nodeIndex)

			if err := file.processOverlays(file.Overlays, nodeIndex); err != nil {
				return fmt.Errorf("failed to apply file overlays, %w", err)
			}

			log.Infof("Processing Document Overlays in File %s on Document %d\n\n", file.Path, nodeIndex)

			for docIndex, doc := range file.Documents {
				if doc.Path != fmt.Sprint(docIndex) {
					continue
				}

				if err := file.processOverlays(file.Documents[docIndex].Overlays, nodeIndex); err != nil {
					return fmt.Errorf("failed to apply document overlays, %w", err)
				}
			}
		}

		if err := file.doPostProcessing(options); err != nil {
			return fmt.Errorf("failed to preform post processing on %s: %w", file.Path, err)
		}
	}

	return nil
}

func (i *Instructions) readYamlFiles() error {
	for index, file := range i.YamlFiles {
		reader, err := ReadStream(file.Path)
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
				return fmt.Errorf("failed to read file %s: %w", file.Path, err)
			}

			i.YamlFiles[index].Nodes = append(i.YamlFiles[index].Nodes, &y)
		}
	}

	return nil
}
