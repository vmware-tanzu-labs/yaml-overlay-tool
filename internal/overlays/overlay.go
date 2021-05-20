// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: MIT

package overlays

import (
	"fmt"

	"github.com/op/go-logging"
	"github.com/vmware-tanzu-labs/yaml-overlay-tool/internal/actions"
	"gopkg.in/yaml.v3"
)

var log = logging.MustGetLogger("overlays") //nolint:gochecknoglobals

type Overlay struct {
	Name          string          `yaml:"name,omitempty"`
	Query         multiString     `yaml:"query,omitempty"`
	Value         yaml.Node       `yaml:"value,omitempty"`
	Action        actions.Action  `yaml:"action,omitempty"`
	DocumentQuery DocumentQueries `yaml:"documentQuery,omitempty"`
	OnMissing     OnMissing       `yaml:"onMissing,omitempty"`
	DocumentIndex []int           `yaml:"documentIndex,omitempty"`
}

func (o *Overlay) Apply(n *yaml.Node) error {
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

func (o *Overlay) CheckDocumentIndex(current int) bool {
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
		case actions.Delete:
			actions.DeleteNode(root, n)
		case actions.Replace:
			if err := actions.ReplaceNode(n, &o.Value); err != nil {
				return fmt.Errorf("%w, skipping replace on line %d column %d", err, n.Line, n.Column)
			}
		case actions.Format:
			if err := actions.FormatNode(n, &o.Value); err != nil {
				return fmt.Errorf("%w, skipping format on line %d column %d", err, n.Line, n.Column)
			}
		case actions.Merge:
			if err := actions.MergeNode(n, &o.Value); err != nil {
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
