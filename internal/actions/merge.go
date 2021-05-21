// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: MIT

package actions

import (
	"errors"
	"fmt"
	"strings"

	"github.com/jinzhu/copier"
	"gopkg.in/yaml.v3"
)

var (
	// ErrMergeMustBeOfSameKind occurs when attempting to merge two different types.
	ErrMergeMustBeOfSameKind = errors.New("both values must be of same type to use 'merge' action")
	// ErrMergeUnsupportedType occurs when an unsupported YAML type is attempted to be merged.
	ErrMergeUnsupportedType = errors.New("not a supported yaml type for merging")
)

func MergeNode(nodes ...*yaml.Node) error {
	if len(nodes) <= 1 {
		return nil
	}

	o := nodes[0]

	for _, n := range nodes[1:] {
		if err := merge(o, n); err != nil {
			return err
		}
	}

	return nil
}

func merge(o, n *yaml.Node, keyName ...string) error {
	switch o.Kind {
	case yaml.DocumentNode:
		return mergeDocument(o, n)
	case yaml.MappingNode:
		return mergeMap(o, n)
	case yaml.SequenceNode:
		return mergeArray(o, n)
	case yaml.ScalarNode:
		mergeScalar(o, n, keyName...)
	case yaml.AliasNode:
		return fmt.Errorf("%s is %w", o.LongTag(), ErrMergeUnsupportedType)
	}

	return nil
}

func mergeDocument(o, n *yaml.Node) error {
	if o.Content != nil && n.Content != nil {
		mergeComments(o, n, o.Value)

		if err := merge(o.Content[0], n.Content[0], o.Value); err != nil {
			return err
		}
	}

	return nil
}

func mergeMap(o, n *yaml.Node) error {
	for ni := 0; ni < len(n.Content)-1; ni += 2 {
		resultFound := false

		for oi := 0; oi < len(o.Content)-1; oi += 2 {
			var formatKey bool

			if formatKey = checkForMarkers(n.Content[ni].Value); !formatKey {
				if o.Content[oi].Value != n.Content[ni].Value {
					continue
				}
			}

			resultFound = true

			keyName := o.Content[oi].Value
			value := o.Content[oi+1].Value

			mergeComments(o.Content[oi], n.Content[ni], value, keyName)

			if err := merge(o.Content[oi+1], n.Content[ni+1], keyName); err != nil {
				return err
			}

			if !formatKey {
				break
			}

			mergeScalar(o.Content[oi], n.Content[ni], value, keyName)
		}

		if !resultFound {
			if err := addNode(o, n.Content[ni:ni+2]...); err != nil {
				return err
			}
		}
	}

	return nil
}

func mergeArray(o, n *yaml.Node) error {
	if o.Content != nil && n.Content != nil {
		mergeComments(o, n, o.Value)

		if err := addNode(o, n.Content...); err != nil {
			return err
		}
	}

	return nil
}

func mergeScalar(o, n *yaml.Node, values ...string) {
	var keyName string

	var value string

	hc := strings.TrimPrefix(o.HeadComment, "#")
	lc := strings.TrimPrefix(o.LineComment, "#")
	fc := strings.TrimPrefix(o.FootComment, "#")

	mergeComments(o, n, values...)

	switch {
	case values == nil:
		value = o.Value
		keyName = ""
	case len(values) > 1:
		value = values[0]
		keyName = values[1]
	case len(values) == 1:
		value = o.Value
		keyName = values[0]
	}

	o.Value = format(n.Value, value, lc, hc, fc, keyName)
}

func mergeComments(o, n *yaml.Node, values ...string) {
	var keyName string

	var value string

	lc := strings.TrimPrefix(o.LineComment, "#")
	hc := strings.TrimPrefix(o.HeadComment, "#")
	fc := strings.TrimPrefix(o.FootComment, "#")

	switch {
	case values == nil:
		value = o.Value
		keyName = ""
	case len(values) > 1:
		value = values[0]
		keyName = values[1]
	case len(values) == 1:
		value = o.Value
		keyName = values[0]
	}

	switch {
	case n.HeadComment != "":
		o.HeadComment = format(n.HeadComment, value, lc, hc, fc, keyName)

		fallthrough
	case n.LineComment != "":
		o.LineComment = format(n.LineComment, value, lc, hc, fc, keyName)

		fallthrough
	case n.FootComment != "":
		o.FootComment = format(n.FootComment, value, lc, hc, fc, keyName)
	}
}

func addNode(o *yaml.Node, nv ...*yaml.Node) error {
	options := copier.Option{
		IgnoreEmpty: false,
		DeepCopy:    true,
	}

	temp := make([]*yaml.Node, len(nv))

	if err := copier.CopyWithOption(&temp, nv, options); err != nil {
		return fmt.Errorf("failed to insert value during merge: %w", err)
	}

	sanitizeNode(temp...)

	o.Content = append(o.Content, temp...)

	return nil
}
