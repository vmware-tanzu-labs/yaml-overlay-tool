// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: MIT

package lib

import (
	"errors"
	"fmt"

	"github.com/vmware-tanzu-labs/yaml-overlay-tool/internal/actions"
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
	DocumentQuery DocumentQueries `yaml:"documentQuery,omitempty"`
	OnMissing     OnMissing       `yaml:"onMissing,omitempty"`
	DocumentIndex []int           `yaml:"documentIndex,omitempty"`
}

func (o *Overlay) apply(n *yaml.Node) error {
	log.Debugf("Checking Document Queries for %s", o.Query)

	if ok, err := o.DocumentQuery.checkQueries(n); !ok {
		return err
	}

	log.Debugf("applying overlay [%q], %s at %s\n", o.Name, o.Action, o.Query)

	results, err := searchYAMLPaths(o.Query, n)
	if err != nil {
		return err
	}

	if results == nil {
		log.Debugf("No results found checking onMissing")

		if err := o.onMissing(n); err != nil {
			return err
		}
	}

	return o.doAction(n, results)
}

func (o *Overlay) checkDocumentIndex(current int) bool {
	if o.DocumentIndex != nil {
		for _, i := range o.DocumentIndex {
			if current != i {
				return true
			}
		}

		return false
	}

	return true
}

func (o *Overlay) doAction(root *yaml.Node, nodes []*yaml.Node) error {
	for _, n := range nodes {
		b, _ := yaml.Marshal(n)
		p, _ := yaml.Marshal(o.Value)

		log.Debugf("Current: >>>\n%s\n", b)
		log.Debugf("Proposed: >>>\n%s\n", p)

		switch o.Action {
		case Delete:
			actions.Delete(root, n)
		case Replace:
			if err := actions.Replace(n, &o.Value); err != nil {
				return fmt.Errorf("%w, skipping replace on line %d column %d", err, n.Line, n.Column)
			}
		case Format:
			if err := actions.Format(n, &o.Value); err != nil {
				return fmt.Errorf("%w, skipping format on line %d column %d", err, n.Line, n.Column)
			}
		case Merge:
			if err := actions.Merge(n, &o.Value); err != nil {
				return fmt.Errorf("%w, skipping merge on line %d, column %d", err, n.Line, n.Column)
			}
		default:
			return fmt.Errorf("%w of type '%s'", ErrInvalidAction, o.Action)
		}
	}

	return nil
}

func (o *Overlay) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type temp Overlay

	if err := unmarshal((*temp)(o)); err != nil {
		return err
	}

	if o.Name == "" {
		o.Name = o.Query.String()
	}

	return nil
}
