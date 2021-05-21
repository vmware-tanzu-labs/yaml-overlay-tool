// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: MIT

package actions

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

func format(formatStr string, v ...interface{}) string {
	formatStr, v = sanitizeMarkers(formatStr, v...)

	v = append(v, "")

	formatStr += fmt.Sprint("%[", len(v), "]s")

	return fmt.Sprintf(formatStr, v...)
}

func sanitizeNode(n ...*yaml.Node) {
	for _, nv := range n {
		if nv == nil {
			continue
		}

		if nv.Kind == yaml.DocumentNode|yaml.MappingNode|yaml.SequenceNode|yaml.AliasNode {
			sanitizeNode(nv.Content...)
		}

		values := []*string{&nv.Value, &nv.HeadComment, &nv.LineComment, &nv.FootComment}
		for _, v := range values {
			*v = format(*v, "", "", "", "", "")
		}
	}
}
