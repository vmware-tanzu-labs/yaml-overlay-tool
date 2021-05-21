// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: MIT

package actions

import (
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

func format(formatStr string, v ...interface{}) string {
	for i, t := range []string{"%v", "%l", "%h", "%f"} {
		for j := range v {
			v[j] = strings.ReplaceAll(v[j].(string), t, "")
		}

		formatStr = strings.ReplaceAll(formatStr, t, fmt.Sprintf("%%[%d]v", i+1))
	}

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
			*v = format(*v, "", "", "", "")
		}
	}
}
