// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: MIT

package actions

import (
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

func format(formatStr string, v ...interface{}) string {
	for i := range v {
		for _, t := range []string{"%v", "%h", "%l", "%f"} {
			v[i] = strings.ReplaceAll(v[i].(string), t, "")
			formatStr = strings.ReplaceAll(formatStr, t, fmt.Sprintf("%%[%d]v", i+1))
		}
	}

	v = append(v, "")

	formatStr += fmt.Sprint("%[", len(v), "]s")

	return fmt.Sprintf(formatStr, v...)
}

func sanatizeNode(n ...*yaml.Node) {
	for _, nv := range n {
		if nv == nil {
			continue
		}

		switch nv.Kind {
		case yaml.DocumentNode, yaml.MappingNode, yaml.SequenceNode, yaml.AliasNode:
			nv.Value = format(nv.Value, "", "", "", "")

			sanatizeNode(nv.Content...)
		case yaml.ScalarNode:
			nv.Value = format(nv.Value, "", "", "", "")
		}
	}
}
