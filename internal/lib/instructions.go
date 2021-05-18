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

func (i *Instructions) addCommonOverlays() {
	for idx := range i.YamlFiles {
		i.YamlFiles[idx].Overlays = append(i.CommonOverlays, i.YamlFiles[idx].Overlays...)
	}
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
