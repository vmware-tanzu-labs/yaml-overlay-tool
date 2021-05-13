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
	Conditions []Condition `yaml:"conditions,omitempty"`
}

type Condition struct {
	Key   string    `yaml:"key,omitempty"`
	Value yaml.Node `yaml:"value,omitempty"`
}

func (dq *DocumentQuery) checkQuery(node *yaml.Node) (bool, error) {
	compareOptions := cmpopts.IgnoreFields(yaml.Node{}, "HeadComment", "LineComment", "FootComment", "Line", "Column", "Style")

	for ci := range dq.Conditions {
		yp, err := yamlpath.NewPath(dq.Conditions[ci].Key)
		if err != nil {
			return false, fmt.Errorf("failed to parse the documentQuery condition %s due to %w", dq.Conditions[ci].Key, err)
		}

		results, err := yp.Find(node)
		if err != nil {
			return false, fmt.Errorf("failed to find results for %s, %w", dq.Conditions[ci].Key, err)
		}

		for _, result := range results {
			if ok := cmp.Equal(*result, dq.Conditions[ci].Value, compareOptions); !ok {
				return false, nil
			}
		}
	}

	return true, nil
}
