// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: MIT

package actions

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

func Format(node, value *yaml.Node) error {
	if node.Kind == yaml.ScalarNode && value.Kind == yaml.ScalarNode {
		value.Value = CondSprintf(value.Value, node.Value)
		*node = *value

		return nil
	}

	return fmt.Errorf("action format can only be used on scalar values")
}

func CondSprintf(format string, v ...interface{}) string {
	v = append(v, "")
	format += fmt.Sprint("%[", len(v), "]s")

	return fmt.Sprintf(format, v...)
}
