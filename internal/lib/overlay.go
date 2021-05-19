// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: MIT

package lib

import (
	"fmt"

	"github.com/vmware-tanzu-labs/yaml-overlay-tool/internal/actions"
	"gopkg.in/yaml.v3"
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
	log.Debugf("Checking Document Queries for [%q]", o.Name)

	if ok, err := o.DocumentQuery.checkQueries(n); !ok {
		return err
	}

	results, err := searchYAMLPaths(o.Query, n)
	if err != nil {
		return err
	}

	if results == nil {
		log.Debugf("No results found checking onMissing")

		return o.onMissing(n)
	}

	return o.doAction(n, results)
}

func (o *Overlay) checkDocumentIndex(current int) bool {
	if o.DocumentIndex != nil {
		for _, i := range o.DocumentIndex {
			if current == i {
				return true
			}
		}

		return false
	}

	return true
}

func (o *Overlay) doAction(root *yaml.Node, nodes []*yaml.Node) error {
	for _, n := range nodes {
		log.Debugf("applying overlay [%q], %s at line %d column %d\n", o.Name, o.Action, n.Line, n.Column)

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
