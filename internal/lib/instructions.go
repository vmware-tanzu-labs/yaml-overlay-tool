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
	CommonOverlays []*Overlay  `yaml:"commonOverlays,omitempty"`
	YamlFiles      []*YamlFile `yaml:"yamlFiles,omitempty"`
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

	instructions.setOutputPath()

	return &instructions, nil
}

func (i *Instructions) addCommonOverlays() {
	for _, yf := range i.YamlFiles {
		yf.Overlays = append(i.CommonOverlays, yf.Overlays...)
	}
}

func (i *Instructions) setOutputPath() {
	p := make([]string, 0, len(i.YamlFiles))

	for _, yf := range i.YamlFiles {
		for _, src := range yf.Sources {
			p = append(p, src.Path)
		}
	}

	pathPrefix := GetCommonPrefix(os.PathSeparator, p...)

	for _, yf := range i.YamlFiles {
		for _, src := range yf.Sources {
			src.outputPath = strings.TrimPrefix(src.Path, pathPrefix)
		}
	}
}
