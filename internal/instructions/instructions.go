// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: MIT

package instructions

import (
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
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

func (cfg *Config) GetInstructions() (*Instructions, error) {
	instructions := new(Instructions)

	if cfg.InstructionsFile != "" {
		instructionsPath, err := filepath.Abs(cfg.InstructionsFile)
		if err != nil {
			return nil, fmt.Errorf("cannot get absolute path of instructions file %s, %w", cfg.InstructionsFile, err)
		}

		instructionsDir = path.Dir(instructionsPath)

		instructions, err = cfg.ReadInstructionFile()
		if err != nil {
			return nil, err
		}
	} else {
		wd, err := os.Getwd()
		if err != nil {
			return nil, fmt.Errorf("could not determine working directory, %w", err)
		}

		instructionsDir = wd
	}

	if err := cfg.ReadAdHocOverlays(instructions); err != nil {
		return nil, err
	}

	instructions.setOutputPath()

	instructions.addCommonOverlays()

	// remove the comments if requested
	if viper.GetBool("removeComments") {
		for _, yamlFile := range instructions.YamlFiles {
			for _, node := range yamlFile.Nodes {
				removeCommentsFromNode(node)
			}
		}
	}

	return instructions, nil
}

// ReadInstructionFile reads a file and decodes it into an Instructions struct.
func (cfg *Config) ReadInstructionFile() (*Instructions, error) {
	var reader io.Reader

	var values interface{}

	log.Debugf("Instructions File: %s\n", cfg.InstructionsFile)

	var instructions Instructions

	instructionsPath, err := filepath.Abs(cfg.InstructionsFile)
	if err != nil {
		return nil, fmt.Errorf("cannot get absolute path of instructions file %s, %w", cfg.InstructionsFile, err)
	}

	instructionsDir = path.Dir(instructionsPath)

	if cfg.Values != nil {
		values, err = getValues(cfg.Values)
		if err != nil {
			return nil, err
		}

		reader, err = renderInstructionsTemplate(cfg.InstructionsFile, values)
		if err != nil {
			return nil, err
		}
	} else {
		reader, err = ReadStream(cfg.InstructionsFile)
		if err != nil {
			return nil, err
		}
	}

	dc := yaml.NewDecoder(reader)
	if err := dc.Decode(&instructions); err != nil {
		return nil, fmt.Errorf("unable to read instructions file %s: %w", cfg.InstructionsFile, err)
	}

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
		if yf.Path != "-" {
			p = append(p, yf.Path)
		}
	}

	pathPrefix := GetCommonPrefix(os.PathSeparator, p...)

	for _, yf := range i.YamlFiles {
		if yf.OutputPath == "" {
			yf.OutputPath = strings.TrimPrefix(yf.Path, pathPrefix)
		} else if path.Ext(yf.OutputPath) == "" {
			yf.OutputPath = path.Join(yf.OutputPath, path.Base(yf.Path))
		}
	}
}
