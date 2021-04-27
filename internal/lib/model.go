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

type Options struct {
	Verbose          bool
	InstructionsFile string
	OutputDir        string
	StdOut           bool
}

type Instructions struct {
	CommonOverlays []Overlay  `yaml:"commonOverlays,omitempty"`
	YamlFiles      []YamlFile `yaml:"yamlFiles,omitempty"`
}

type Overlay struct {
	Name          string          `yaml:"name,omitempty"`
	Query         string          `yaml:"query,omitempty"`
	Value         yaml.Node       `yaml:"value,omitempty"`
	Action        string          `yaml:"action,omitempty"`
	DocumentQuery []DocumentQuery `yaml:"documentQuery,omitempty"`
	OnMissing     OnMissing       `yaml:"onMissing,omitempty"`
	DocumentIndex []int           `yaml:"documentIndex,omitempty"`
}

type DocumentQuery struct {
	Conditions []Condition `yaml:"conditions,omitempty"`
}

type Condition struct {
	Key   string `yaml:"key,omitempty"`
	Value string `yaml:"value,omitempty"`
}

type YamlFile struct {
	Name      string     `yaml:"name,omitempty"`
	Path      string     `yaml:"path,omitempty"`
	Overlays  []Overlay  `yaml:"overlays,omitempty"`
	Documents []YamlFile `yaml:"documents,omitempty"`
	Nodes     []*yaml.Node
}

type OnMissing struct {
	Action     string      `yaml:"action,omitempty"`
	InjectPath interface{} `yaml:"injectPath,omitempty"`
}

type YamlDocument struct {
	Name     string    `yaml:"name,omitempty"`
	Index    int       `yaml:"path,omitempty"`
	Overlays []Overlay `yaml:"overlays,omitempty"`
}

func (i *Instructions) ReadYamlFiles() error {
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
