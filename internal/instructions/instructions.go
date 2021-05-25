// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: MIT

package instructions

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/vmware-tanzu-labs/yaml-overlay-tool/internal/overlays"
	"gopkg.in/yaml.v3"
)

var instructionsDir string //nolint:gochecknoglobals // needed for path lookups relative to instructions file

// Instructions is a struct used for decoding an instructions file.
type Instructions struct {
	// Common Overlays that will apply to all files specified.
	CommonOverlays []*overlays.Overlay `yaml:"commonOverlays,omitempty"`
	// List of YamlFiles and overlays to apply.
	YamlFiles YamlFiles `yaml:"yamlFiles,omitempty"`
}

// ReadInstructionFile reads a file and decodes it into an Instructions struct.
func ReadInstructionFile(fileName *string) (*Instructions, error) {
	var instructions Instructions

	var err error

	log.Debugf("Instructions File: %s\n", *fileName)

	instructionsPath, err := filepath.Abs(*fileName)
	if err != nil {
		return nil, fmt.Errorf("cannot get absolute path of instructions file %s, %w", *fileName, err)
	}

	instructionsDir = path.Dir(instructionsPath)

	reader, err := ReadStream(instructionsPath)
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

// addCommonOverlays, takes all common overlays and adds them to each yamlFile entry.
func (i *Instructions) addCommonOverlays() {
	for _, yf := range i.YamlFiles {
		yf.Overlays = append(i.CommonOverlays, yf.Overlays...)
	}
}

func (i *Instructions) setOutputPath() {
	p := make([]string, 0, len(i.YamlFiles))

	for _, yf := range i.YamlFiles {
		p = append(p, yf.Path)
	}

	pathPrefix := GetCommonPrefix(os.PathSeparator, p...)

	for _, yf := range i.YamlFiles {
		yf.OutputPath = strings.TrimPrefix(yf.Path, pathPrefix)
	}
}
