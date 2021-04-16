// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: MIT

package edit

import (
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

func IterateNode(node *yaml.Node, identifier string) *yaml.Node {
	returnNode := false

	for _, n := range node.Content {
		if n.Value == identifier {
			returnNode = true
			continue
		}

		if returnNode {
			return n
		}

		if len(n.Content) > 0 {
			acNode := IterateNode(n, identifier)
			if acNode != nil {
				return acNode
			}
		}
	}

	return nil
}

func IteratePath(node *yaml.Node, path string) (*yaml.Node, error) {
	components := strings.Split(path, ".")
	pn := node

	for _, c := range components {
		node = IterateNode(pn, c)
		if node == nil {
			return pn, fmt.Errorf("subpath: %s not found in Path: %s ", c, path)
		}

		pn = node
	}

	return node, nil
}
