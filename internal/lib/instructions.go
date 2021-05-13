// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: MIT

package lib

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

type Instructions struct {
	CommonOverlays []Overlay  `yaml:"commonOverlays,omitempty"`
	YamlFiles      []YamlFile `yaml:"yamlFiles,omitempty"`
}

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

	if err := instructions.ReadYamlFiles(); err != nil {
		return nil, err
	}

	instructions.setOutputPath()

	return &instructions, nil
}

func (i *Instructions) processYamlFiles(options *Options) error {
	for _, file := range i.YamlFiles {
		for _, src := range file.Source {
			for nodeIndex := range src.Nodes {
				log.Infof("Processing Common Overlays in File %s on Document %d\n\n", src.Path, nodeIndex)

				if err := src.processOverlays(i.CommonOverlays, nodeIndex); err != nil {
					return fmt.Errorf("failed to apply common overlays, %w", err)
				}

				log.Infof("Processing File Overlays in File %s on Document %d\n\n", src.Path, nodeIndex)

				if err := src.processOverlays(file.Overlays, nodeIndex); err != nil {
					return fmt.Errorf("failed to apply file overlays, %w", err)
				}

				log.Infof("Processing Document Overlays in File %s on Document %d\n\n", src.Path, nodeIndex)

				for di := range file.Documents {
					if !checkDocumentPath(&file.Documents[di], nodeIndex) {
						continue
					}

					if err := src.processOverlays(file.Documents[di].Overlays, nodeIndex); err != nil {
						return err
					}
				}
			}

			if err := src.doPostProcessing(options); err != nil {
				return fmt.Errorf("failed to perform post processing on %s: %w", src.Path, err)
			}
		}
	}

	return nil
}

func (i *Instructions) ReadYamlFiles() error {
	for yIndex, yFile := range i.YamlFiles {
		var files []string

		if ok, err := isDirectory(yFile.Path); ok {
			files, err = getFileNames(yFile.Path)
			if err != nil {
				return err
			}
		} else {
			if err != nil {
				return err
			}

			files = []string{yFile.Path}
		}

		for _, file := range files {
			if err := i.YamlFiles[yIndex].readYamlFile(file); err != nil {
				return err
			}
		}
	}

	return nil
}

func (i *Instructions) setOutputPath() {
	p := make([]string, 0, len(i.YamlFiles))

	for _, yf := range i.YamlFiles {
		for _, src := range yf.Source {
			p = append(p, src.Path)
		}
	}

	pathPrefix := GetCommonPrefix(os.PathSeparator, p...)

	for yi := range i.YamlFiles {
		for si, src := range i.YamlFiles[yi].Source {
			i.YamlFiles[yi].Source[si].outputPath = strings.TrimPrefix(src.Path, pathPrefix)
		}
	}
}

func checkDocumentPath(doc *YamlFile, docIndex int) bool {
	return doc.Path == fmt.Sprint(docIndex)
}
