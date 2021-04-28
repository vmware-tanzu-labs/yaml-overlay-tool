// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: MIT

package lib

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/vmware-labs/yaml-jsonpath/pkg/yamlpath"
	"gopkg.in/yaml.v3"
)

type Options struct {
	Verbose          bool
	InstructionsFile string
	OutputDir        string
	StdOut           bool
	Indent           int
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

func (f *YamlFile) processOverlays(o []Overlay, nodeIndex int) error {
	for i := range o {
		if err := o[i].process(f, nodeIndex); err != nil {
			return err
		}
	}

	return nil
}

func (o *Overlay) process(f *YamlFile, i int) error {
	if ok := checkDocIndex(i, o.DocumentIndex); !ok {
		return nil
	}

	node := f.Nodes[i]

	ok, err := o.doDocumentQuery(node)
	if err != nil {
		return err
	}

	if !ok {
		return nil
	}

	log.Debugf("%s at %s in file %s on Document %d\n", o.Action, o.Query, f.Path, i)

	yp, err := yamlpath.NewPath(o.Query)
	if err != nil {
		return fmt.Errorf("failed to parse the query path %s due to %w", o.Query, err)
	}

	results, err := yp.Find(node)
	if err != nil {
		return fmt.Errorf("faild to find results for %s, %w", o.Query, err)
	}

	if results == nil {
		log.Debugf("Call OnMissing Here")
		o.onMissing(f, i)
	}

	if err := processActions(node, results, o); err != nil {
		if errors.Is(err, ErrInvalidAction) {
			return fmt.Errorf("%w in instructions file", err)
		}

		return fmt.Errorf("%w in file %s on Document %d", err, f.Path, i)
	}

	return nil
}

func (f *YamlFile) Save(o *Options) error {
	b, err := yaml.Marshal(f.Nodes)
	if err != nil {
		return fmt.Errorf("failed to marshal %s: %w", f.Path, err)
	}

	fileName := path.Base(f.Path)
	outputFileName := path.Join(o.OutputDir, fileName)
	//nolint:gosec //output files with read and write permissions so that end-users can continue to leverage these files
	if err := ioutil.WriteFile(outputFileName, b, 0644); err != nil {
		return fmt.Errorf("failed to save file %s: %w", outputFileName, err)
	}

	return nil
}

func (f *YamlFile) doPostProcessing(o *Options) error {
	output := new(bytes.Buffer)
	ye := yaml.NewEncoder(output)
	ye.SetIndent(o.Indent)

	for _, node := range f.Nodes {
		err := ye.Encode(node)
		if err != nil {
			return fmt.Errorf("unable to marshal final document %s, error: %w", f.Path, err)
		}
	}

	log.Noticef("Final: >>>\n%s\n", output)
	// added so we can quickly see the results of the run
	if o.StdOut {
		fmt.Printf("---\n%s", output) //nolint:forbidigo
	}

	return nil
}

func (o *Overlay) doDocumentQuery(node *yaml.Node) (bool, error) {
	log.Debugf("Checking Document Queries for %s", o.Query)

	if o.DocumentQuery == nil {
		log.Debugf("No Document Queries found")

		return true, nil
	}

	conditionsMet := false

	options := cmpopts.IgnoreFields(yaml.Node{}, "HeadComment", "LineComment", "FootComment", "Line", "Column", "Style")

	for i := range o.DocumentQuery {
		for ci := range o.DocumentQuery[i].Conditions {
			yp, err := yamlpath.NewPath(o.DocumentQuery[i].Conditions[ci].Key)
			if err != nil {
				return false, fmt.Errorf("failed to parse the query path %s due to %w", o.DocumentQuery[i].Conditions[ci].Key, err)
			}

			results, err := yp.Find(node)
			if err != nil {
				return false, fmt.Errorf("failed to find results for %s, %w", o.DocumentQuery[i].Conditions[ci].Key, err)
			}

			for _, result := range results {
				conditionsMet = cmp.Equal(*result, o.DocumentQuery[i].Conditions[ci].Value, options)
				if !conditionsMet {
					break
				}
			}
		}

		if conditionsMet {
			log.Debugf("Document Query conditions were met, continuing")

			return true, nil
		}
	}

	log.Debugf("Document Query Conditions were not met, skipping")

	return false, nil
}
