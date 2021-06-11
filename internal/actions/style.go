// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: MIT

package actions

import (
	"strings"

	"gopkg.in/yaml.v3"
)

type Style uint

type Styles []Style

const (
	// TaggedStyle chooses YAML Tag style for rendering.
	TaggedStyle Style = 1 << iota
	// DoubleQuotedStyle chooses YAML Double Quoted style for rendering.
	DoubleQuotedStyle
	// SingleQuotedStyle chooses YAML Single Quoted style for rendering.
	SingleQuotedStyle
	// LiteralStyle chooses use YAML Literal style for rendering.
	LiteralStyle
	// FoldedStyle chooses use YAML Folded style for rendering.
	FoldedStyle
	// FlowStyle chooses YAML Flow style for rendering.
	FlowStyle
	// NormalStyle chooses YAML Normal style for rendering.
	NormalStyle
)

func (s Style) String() string {
	toString := map[Style]string{
		NormalStyle:       "normal",
		TaggedStyle:       "tagged",
		DoubleQuotedStyle: "doubleQuoted",
		SingleQuotedStyle: "singleQuoted",
		LiteralStyle:      "literal",
		FoldedStyle:       "folded",
		FlowStyle:         "flow",
	}

	return toString[s]
}

func (ss Styles) FlagMap() map[Style][]string {
	return map[Style][]string{
		NormalStyle:       {"normal", "n"},
		TaggedStyle:       {"tagged", "tag", "t"},
		DoubleQuotedStyle: {"doubleQuoted", "doubleQuote", "double", "dq"},
		SingleQuotedStyle: {"singleQuoted", "singleQoute", "single", "sq"},
		LiteralStyle:      {"literal", "l"},
		FoldedStyle:       {"folded", "fold", "fo"},
		FlowStyle:         {"flow", "fl"},
	}
}

func (ss Styles) GetStyle() Style {
	var s Style
	for _, v := range ss {
		s += v
	}

	return s
}

func SetStyle(ss Styles, n ...*yaml.Node) {
	s := ss.GetStyle()
	setStyle(s, n...)
}

func setStyle(s Style, n ...*yaml.Node) {
	for _, nv := range n {
		switch nv.Kind {
		case yaml.DocumentNode, yaml.MappingNode, yaml.SequenceNode, yaml.AliasNode:
			nv.Style = yaml.Style(s)

			setStyle(s, nv.Content...)
		case yaml.ScalarNode:
			nv.Style = yaml.Style(s)
		}
	}
}

func GetStyleFromConfig(v ...string) Styles {
	var result Styles

	for _, ss := range v {
		ssp := strings.Split(strings.Trim(strings.ToLower(ss), "[]"), ",")
		for _, s := range ssp {
			switch s {
			case "normal", "n":
				result = append(result, NormalStyle)
			case "tagged", "tag", "t":
				result = append(result, TaggedStyle)
			case "doublequoted", "doublequote", "double", "dq":
				result = append(result, DoubleQuotedStyle)
			case "singlequoted", "singlequote", "single", "sq":
				result = append(result, SingleQuotedStyle)
			case "literal", "l":
				result = append(result, LiteralStyle)
			case "folded", "fold", "fo":
				result = append(result, FoldedStyle)
			case "flow", "fl":
				result = append(result, FlowStyle)
			}
		}
	}

	return result
}
