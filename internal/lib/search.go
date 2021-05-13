// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: MIT

package lib

import (
	"fmt"

	"github.com/vmware-labs/yaml-jsonpath/pkg/yamlpath"
	"gopkg.in/yaml.v3"
)

func searchYAMLPaths(paths []string, node *yaml.Node) ([]*yaml.Node, error) {
	var results []*yaml.Node

	for _, p := range paths {
		log.Debugf("searching path %s\n", p)

		yp, err := yamlpath.NewPath(p)
		if err != nil {
			return nil, fmt.Errorf("failed to parse the query path %s due to %w", p, err)
		}

		result, err := yp.Find(node)
		if err != nil {
			return nil, fmt.Errorf("failed to find results for %s, %w", p, err)
		}

		results = append(results, result...)
	}

	return results, nil
}
