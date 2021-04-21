// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: MIT

package actions

import "gopkg.in/yaml.v3"

func DeleteNode(n *yaml.Node) {
	n.Content = []*yaml.Node{}
}

func DeleteSeqNode(n *yaml.Node, key, value string) {
	state := -1
	indexRemove := -1

	for index, pn := range n.Content {
		for _, cn := range pn.Content {
			if key == cn.Value && state == -1 {
				state++

				continue // found expected move onto next
			}

			if value == cn.Value && state == 0 {
				state++

				indexRemove = index

				break // found the target exit out of the loop
			} else if state == 0 {
				state = -1
			}
		}
	}

	if state == 1 {
		// Remove node from contents
		length := len(n.Content)
		copy(n.Content[indexRemove:], n.Content[indexRemove+1:])
		n.Content[length-1] = nil
		n.Content = n.Content[:length-1]
	}
}
