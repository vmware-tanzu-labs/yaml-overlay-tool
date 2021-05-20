// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: MIT

package actions

import (
	"errors"
	"fmt"
	"strconv"

	"gopkg.in/yaml.v3"
)

var (
	ErrMergeMustBeOfSameKind = errors.New("both values must be of same type to use 'merge' action")
	ErrMergeUnsupportedType  = errors.New("not a supported yaml type for merging")
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

func merge(o, n *yaml.Node) error {
	if o.Kind != n.Kind && o.Kind != 0 {
		// are both originalValue and newValue the same 'Kind'?
		return ErrMergeMustBeOfSameKind
	}

	switch o.Kind {
	case yaml.DocumentNode:
		return mergeDocument(o, n)
	case yaml.MappingNode:
		return mergeMap(o, n)
	case yaml.SequenceNode:
		mergeArray(o, n)
	case yaml.ScalarNode:
		return mergeScalar(o, n)
	case yaml.AliasNode:
		return fmt.Errorf("%s is %w", o.LongTag(), ErrMergeUnsupportedType)
	}

	return nil
}

func mergeDocument(o, n *yaml.Node) error {
	if o.Content != nil && n.Content != nil {
		if err := merge(o.Content[0], n.Content[0]); err != nil {
			return err
		}

		mergeComments(o, n)
	}

	return nil
}

func mergeMap(o, n *yaml.Node) error {
	if o.Content != nil && n.Content != nil {
		for ni := 0; ni < len(n.Content)-1; ni += 2 {
			resultFound := false

			for oi := 0; oi < len(o.Content)-1; oi += 2 {
				if o.Content[oi].Value == n.Content[ni].Value {
					resultFound = true

					if err := merge(o.Content[oi+1], n.Content[ni+1]); err != nil {
						return err
					}

					mergeComments(o.Content[oi], n.Content[ni])

					break
				}
			}

			if !resultFound {
				o.Content = append(o.Content, n.Content[ni:ni+2]...)
			}
		}
	}

	return nil
}

func mergeArray(o, n *yaml.Node) {
	if o.Content != nil && n.Content != nil {
		o.Content = append(o.Content, n.Content...)
	}
}

func mergeScalar(ov, nv *yaml.Node) error {
	switch ov.Tag {
	case "!!int":
		o, err := strconv.Atoi(ov.Value)
		if err != nil {
			return fmt.Errorf("failed to %s to int for merging: %w", ov.Value, err)
		}

		n, err := strconv.Atoi(nv.Value)
		if err != nil {
			return fmt.Errorf("failed to %s to int for merging: %w", nv.Value, err)
		}

		ov.Value = strconv.Itoa(o + n)

	case "!!bool":
		o, err := strconv.ParseBool(ov.Value)
		if err != nil {
			return fmt.Errorf("failed to %s to bool for merging: %w", ov.Value, err)
		}

		n, err := strconv.ParseBool(nv.Value)
		if err != nil {
			return fmt.Errorf("failed to %s to bool for merging: %w", nv.Value, err)
		}

		ov.Value = strconv.FormatBool(o && n)

	default:
		ov.Value += nv.Value
	}

	mergeComments(ov, nv)

	return nil
}

func mergeComments(o, n *yaml.Node) {
	switch {
	case n.HeadComment != "":
		o.HeadComment += n.HeadComment

		fallthrough
	case n.LineComment != "":
		o.LineComment += n.LineComment

		fallthrough
	case n.FootComment != "":
		o.FootComment += n.FootComment
	}
}
