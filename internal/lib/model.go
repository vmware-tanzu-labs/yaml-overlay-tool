// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: MIT

package lib

import (
	"gopkg.in/yaml.v3"
)

type Options struct {
	Verbose          bool
	InstructionsFile string
	OutputDir        string
	StdOut           bool
	Indent           int
	ValuesPath       []string
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
	Key   string    `yaml:"key,omitempty"`
	Value yaml.Node `yaml:"value,omitempty"`
}

type YamlFile struct {
	Name       string     `yaml:"name,omitempty"`
	Path       string     `yaml:"path,omitempty"`
	Overlays   []Overlay  `yaml:"overlays,omitempty"`
	Documents  []YamlFile `yaml:"documents,omitempty"`
	Nodes      []*yaml.Node
	outputPath string
}

type OnMissing struct {
	Action     string      `yaml:"action,omitempty"`
	InjectPath interface{} `yaml:"injectPath,omitempty"`
}
