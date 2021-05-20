// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: MIT

package actions

import (
	"errors"
	"fmt"
	"strings"

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
		mergeComments(o, n)

		if err := merge(o.Content[0], n.Content[0]); err != nil {
			return err
		}
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

					mergeComments(o.Content[oi], n.Content[ni])

					if err := merge(o.Content[oi+1], n.Content[ni+1]); err != nil {
						return err
					}

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
		mergeComments(o, n)

		o.Content = append(o.Content, n.Content...)
	}
}

func mergeScalar(ov, nv *yaml.Node) error {
	mergeComments(ov, nv)

	ov.Value = CondSprintf(nv.Value, ov.Value, strings.TrimPrefix(ov.LineComment, "#"))

	return nil
}

func mergeComments(o, n *yaml.Node) {
	hc := strings.TrimPrefix(o.HeadComment, "#")
	lc := strings.TrimPrefix(o.LineComment, "#")
	fc := strings.TrimPrefix(o.FootComment, "#")

	switch {
	case n.HeadComment != "":
		o.HeadComment = CondSprintf(n.HeadComment, o.Value, lc, hc, fc)

		fallthrough
	case n.LineComment != "":
		o.LineComment = CondSprintf(n.LineComment, o.Value, lc, hc, fc)

		fallthrough
	case n.FootComment != "":
		o.FootComment = CondSprintf(n.FootComment, o.Value, lc, hc, fc)
	}
}
