// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: MIT

package overlays

import (
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"gopkg.in/yaml.v3"
)

type DocumentQuery struct {
	Conditions []*Condition `yaml:"conditions,omitempty"`
}

type Condition struct {
	Query Queries   `yaml:"query,omitempty"`
	Value yaml.Node `yaml:"value,omitempty"`
}

type DocumentQueries []*DocumentQuery

func (dq DocumentQueries) checkQueries(node *yaml.Node) bool {
	if len(dq) == 0 {
		log.Debugf("No Document Queries found, continuing")

		return true
	}

	for _, q := range dq {
		if ok := q.checkQuery(node); ok {
			log.Debugf("Document Query conditions were met, continuing")

			return true
		}
	}

	log.Debugf("Document Query conditions were not met, skipping")

	return false
}

func (dq *DocumentQuery) checkQuery(node *yaml.Node) bool {
	compareOptions := cmpopts.IgnoreFields(yaml.Node{}, "HeadComment", "LineComment", "FootComment", "Line", "Column", "Style")

	for _, c := range dq.Conditions {
		results := c.Query.Find(node)

		for _, result := range results {
			if ok := cmp.Equal(*result, c.Value, compareOptions); !ok {
				return false
			}
		}
	}

	return true
}
