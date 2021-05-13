// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: MIT

package lib

import (
	"errors"
	"fmt"
	"strings"

	"github.com/vmware-tanzu-labs/yaml-overlay-tool/internal/actions"
	"github.com/vmware-tanzu-labs/yaml-overlay-tool/internal/path"
	"gopkg.in/yaml.v3"
)

var (
	ErrOnMissingNoInjectAction = errors.New("no matches and no onMissing.action of 'inject'")
	ErrOnMissingNoInjectPath   = errors.New("no matches and no onMissing.injectPath")
	ErrOnMissingInvalidType    = errors.New("invalid type for onMissing.injectPath")
)

type Overlay struct {
	Name          string          `yaml:"name,omitempty"`
	Query         multiString     `yaml:"query,omitempty"`
	Value         yaml.Node       `yaml:"value,omitempty"`
	Action        Action          `yaml:"action,omitempty"`
	DocumentQuery []DocumentQuery `yaml:"documentQuery,omitempty"`
	OnMissing     OnMissing       `yaml:"onMissing,omitempty"`
	DocumentIndex []int           `yaml:"documentIndex,omitempty"`
}

type OnMissing struct {
	Action     OnMissingAction `yaml:"action,omitempty"`
	InjectPath multiString     `yaml:"injectPath,omitempty"`
}

type multiString []string

func (ms *multiString) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var s string
	if err := unmarshal(&s); err == nil {
		*ms = []string{s}

		return nil
	}

	type ss []string

	return unmarshal((*ss)(ms))
}

func (ms multiString) String() string {
	return strings.Join(ms, ",")
}

func (o *Overlay) applyOverlay(src *Source, docIndex int) error {
	if ok := o.checkDocumentIndex(docIndex); !ok {
		return nil
	}

	node := src.Nodes[docIndex]

	if ok, err := o.checkDocumentQuery(node); !ok {
		return err
	}

	log.Debugf("%s at %s in file %s on Document %d\n", o.Action, o.Query, src.Path, docIndex)

	results, err := searchYAMLPaths(o.Query, node)
	if err != nil {
		return err
	}

	if results == nil {
		if err := o.onMissing(src, docIndex); err != nil {
			return err
		}
	}

	if err := o.doAction(node, results); err != nil {
		if errors.Is(err, ErrInvalidAction) {
			return fmt.Errorf("%w in instructions file", err)
		}

		return fmt.Errorf("%w in file %s on Document %d", err, src.Path, docIndex)
	}

	return nil
}

func (o *Overlay) checkDocumentIndex(current int) bool {
	if o.DocumentIndex != nil {
		for f := range o.DocumentIndex {
			if current == o.DocumentIndex[f] {
				return true
			}
		}

		return false
	}

	return true
}

func (o *Overlay) checkDocumentQuery(node *yaml.Node) (bool, error) {
	log.Debugf("Checking Document Queries for %s", o.Query)

	if o.DocumentQuery == nil {
		log.Debugf("No Document Queries found, continuing")

		return true, nil
	}

	for i := range o.DocumentQuery {
		if ok, err := o.DocumentQuery[i].checkQuery(node); ok {
			log.Debugf("Document Query conditions were met, continuing")

			return true, nil
		} else if err != nil {
			return false, err
		}
	}

	log.Debugf("Document Query Conditions were not met, skipping")

	return false, nil
}

func (o *Overlay) onMissing(src *Source, docIndex int) error {
	// check if the query has a match
	// if no match then we require an inject path
	// we need to then check if each inject path is valid (does it exist)
	// if we had an inject path(s) then we inject the value to those locations
	// if we didn't have an inject path we have an implicit onMissing: ignore and we put out a warning if not stdout option to terminal
	switch o.OnMissing.Action {
	case Ignore:
		log.Debugf("ignoring %s at %s in file %s on Document %d due to %s\n", o.Action, o.Query, src.Path, docIndex, ErrOnMissingNoInjectAction)

		return nil
	case Inject:
		_, err := path.BuildMulti(o.Query)
		if err != nil {
			if errors.Is(err, path.ErrInvalidPathSyntax) {
				return o.handleInjectPath(src, docIndex)
			}

			return fmt.Errorf("%w, for onMissing", err)
		}

		return o.doInjectPath(o.Query, src.Nodes[docIndex])
	default:
		return fmt.Errorf("%w for onMissing of type '%s'", ErrInvalidAction, o.Action)
	}
}

func (o *Overlay) doInjectPath(ip []string, node *yaml.Node) error {
	y, err := path.BuildMulti(ip)
	if err != nil {
		return fmt.Errorf("failed to build inject path %s, %w", ip, err)
	}

	err = actions.Merge(node, y)
	if err != nil {
		return fmt.Errorf("failed to merge injectpath scaffolding %s with document, %w", ip, err)
	}

	results, err := searchYAMLPaths(ip, node)
	if err != nil {
		return fmt.Errorf("%w, on injectPath %s", err, ip)
	}

	for i := range results {
		if err := actions.Replace(results[i], &o.Value); err != nil {
			if errors.Is(err, ErrInvalidAction) {
				return fmt.Errorf("%w in instructions file", err)
			}

			return fmt.Errorf("%w for onMissing.InjectPath", err)
		}
	}

	return nil
}

func (o *Overlay) doAction(root *yaml.Node, nodes []*yaml.Node) error {
	for i := range nodes {
		b, _ := yaml.Marshal(nodes[i])
		p, _ := yaml.Marshal(o.Value)

		log.Debugf("Current: >>>\n%s\n", b)
		log.Debugf("Proposed: >>>\n%s\n", p)

		switch o.Action {
		case Delete:
			actions.Delete(root, nodes[i])
		case Replace:
			if err := actions.Replace(nodes[i], &o.Value); err != nil {
				return fmt.Errorf("%w, skipping replace", err)
			}
		case Format:
			if err := actions.Format(nodes[i], &o.Value); err != nil {
				return fmt.Errorf("%w, skipping format", err)
			}
		case Merge:
			if err := actions.Merge(nodes[i], &o.Value); err != nil {
				return fmt.Errorf("%w, skipping merge", err)
			}
		default:
			return fmt.Errorf("%w of type '%s'", ErrInvalidAction, o.Action)
		}
	}

	return nil
}

func (o *Overlay) handleInjectPath(src *Source, docIndex int) error {
	_, err := path.BuildMulti(o.Query)
	if !errors.Is(err, path.ErrInvalidPathSyntax) {
		return o.doInjectPath(o.Query, src.Nodes[docIndex])
	}

	if o.OnMissing.InjectPath == nil {
		log.Debugf("ignoring %s at %s in file %s on Document %d due to %s\n", o.Action, o.Query, src.Path, docIndex, ErrOnMissingNoInjectPath)

		return nil
	}

	if err := o.doInjectPath(o.OnMissing.InjectPath, src.Nodes[docIndex]); err != nil {
		return fmt.Errorf("%w in file %s on Document %d", err, src.Path, docIndex)
	}

	return nil
}
