// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: MIT

package actions

import (
	"strings"

	"github.com/vmware-labs/yaml-jsonpath/pkg/yamlpath"
	"gopkg.in/yaml.v3"
)

// TODO: add delete caps for sequence nodes.
func Delete(root, child *yaml.Node, path string) error {
	pc := strings.Split(path, ".")
	parentPath := strings.Join(pc[:len(pc)-1], ".")

	// if we are searching at root we need to add the root anchor to unwrap the document node so thing process correctly
	if parentPath == "" {
		parentPath = "$"
	}

	yp, err := yamlpath.NewPath(parentPath)
	if err != nil {
		return err
	}

	parentNodes, err := yp.Find(root)
	if err != nil {
		return err
	}

	for _, pn := range parentNodes {
		if err := deleteNode(pn, child); err != nil {
			return err
		}
	}

	return nil
}

func deleteNode(pn, child *yaml.Node) error {
	for i, c := range pn.Content {
		if c != nil {
			if c.Kind == yaml.SequenceNode {
				if err := deleteNode(c, child); err != nil {
					return err
				}
			}
		}

		if c != child {
			continue
		}

		length := len(pn.Content)
		floor := 1
		roof := 2

		if pn.Kind == yaml.SequenceNode {
			floor--
			roof--
		}

		copy(pn.Content[i-floor:], pn.Content[i+1:])
		pn.Content[length-roof] = nil
		pn.Content = pn.Content[:length-roof]

		return nil
	}

	return nil
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
