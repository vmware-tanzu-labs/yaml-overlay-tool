// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: MIT

package instructions

import (
	"fmt"

	"github.com/vmware-tanzu-labs/yaml-overlay-tool/internal/actions"
	"gopkg.in/yaml.v3"
)

func getValues(fileNames []string) (interface{}, error) {
	var values interface{}

	yamlValues := make([]*yaml.Node, len(fileNames))

	for i, v := range fileNames {
		var yamlValue yaml.Node

		reader, err := ReadStream(v)
		if err != nil {
			return nil, err
		}

		yd := yaml.NewDecoder(reader)

		if err := yd.Decode(&yamlValue); err != nil {
			return nil, fmt.Errorf("failed to decode yaml values: %w", err)
		}

		yamlValues[i] = &yamlValue
	}

	if err := actions.MergeNode(yamlValues...); err != nil {
		return nil, fmt.Errorf("failed to merge yaml values: %w", err)
	}

	b, err := yaml.Marshal(yamlValues[0])
	if err != nil {
		return nil, fmt.Errorf("failed to marshal yaml values: %w", err)
	}

	if err := yaml.Unmarshal(b, &values); err != nil {
		return nil, fmt.Errorf("failed to unmarshal yaml values: %w", err)
	}

	return values, nil
}
