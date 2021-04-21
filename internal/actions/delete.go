// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: MIT

package actions

import (
	"gopkg.in/yaml.v3"
)

func Delete(pn, child *yaml.Node) {
	for i, c := range pn.Content {
		// we are comparing against the memory address of the pointer not the value
		// so this will only find one result in the yaml tree
		if c != child {
			if c.Content != nil {
				Delete(c, child)
			}

			continue
		}

		// when working with sequance nodes we only need to delete one node
		// however when working with scalars or maps we need to delete both the key and the value
		length := len(pn.Content)

		// start is the the subtractor needed to get to the first key to delete
		// for maps/scalars this would be -1 since the key is always the index right before the value
		// for sequences this would be 0 since we only have one vlaue to delete, the one we are currently on
		start := 1

		// nodes to delete is the amount of nodes needed to be deleted
		// for sequence nodes this would be 1
		// for maps and scalars this would be 2
		nodesToDelete := 2

		if pn.Kind == yaml.SequenceNode {
			start--
			nodesToDelete--
		}

		copy(pn.Content[i-start:], pn.Content[i+1:])
		pn.Content[length-nodesToDelete] = nil
		pn.Content = pn.Content[:length-nodesToDelete]

		return
	}
}
