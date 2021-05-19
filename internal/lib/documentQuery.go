// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: MIT

package lib

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
	Key   string    `yaml:"key,omitempty"`
	Value yaml.Node `yaml:"value,omitempty"`
}

type DocumentQueries []*DocumentQuery

func (dq DocumentQueries) checkQueries(node *yaml.Node) (bool, error) {
	if dq == nil {
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

	log.Debugf("Document Query Conditions were not met, skipping")

	return false, nil
}

func (dq *DocumentQuery) checkQuery(node *yaml.Node) (bool, error) {
	compareOptions := cmpopts.IgnoreFields(yaml.Node{}, "HeadComment", "LineComment", "FootComment", "Line", "Column", "Style")

	for _, c := range dq.Conditions {
		yp, err := yamlpath.NewPath(c.Key)
		if err != nil {
			return false, fmt.Errorf("failed to parse the documentQuery condition %s due to %w", c.Key, err)
		}

		results, err := yp.Find(node)
		if err != nil {
			return false, fmt.Errorf("failed to find results for %s, %w", c.Key, err)
		}

		for _, result := range results {
			if ok := cmp.Equal(*result, c.Value, compareOptions); !ok {
				return false, nil
			}
		}
	}

	return true, nil
}
