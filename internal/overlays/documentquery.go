// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: MIT

package overlays

import (
	"fmt"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/vmware-labs/yaml-jsonpath/pkg/yamlpath"
	"gopkg.in/yaml.v3"
)

type DocumentQuery struct {
	Conditions []*Condition `yaml:"conditions,omitempty"`
}

type Condition struct {
	Query string    `yaml:"query,omitempty"`
	Value yaml.Node `yaml:"value,omitempty"`
}

type DocumentQueries []*DocumentQuery

func (dq DocumentQueries) checkQueries(node *yaml.Node) (bool, error) {
	if len(dq) == 0 {
		log.Debugf("No Document Queries found, continuing")

		return true, nil
	}

	for _, q := range dq {
		if ok, err := q.checkQuery(node); ok {
			log.Debugf("Document Query conditions were met, continuing")

			return true, nil
		} else if err != nil {
			return false, err
		}
	}

	log.Debugf("Document Query conditions were not met, skipping")

	return false, nil
}

func (dq *DocumentQuery) checkQuery(node *yaml.Node) (bool, error) {
	compareOptions := cmpopts.IgnoreFields(yaml.Node{}, "HeadComment", "LineComment", "FootComment", "Line", "Column", "Style")

	for _, c := range dq.Conditions {
		yp, err := yamlpath.NewPath(c.Query)
		if err != nil {
			return false, fmt.Errorf("failed to parse the documentQuery condition %s due to %w", c.Query, err)
		}

		results, _ := yp.Find(node)

		for _, result := range results {
			if ok := cmp.Equal(*result, c.Value, compareOptions); !ok {
				return false, nil
			}
		}
	}

	return true, nil
}
