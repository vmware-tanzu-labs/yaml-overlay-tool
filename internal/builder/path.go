// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: MIT

package builder

import (
	"errors"
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/vmware-tanzu-labs/yaml-overlay-tool/internal/actions"
	"gopkg.in/yaml.v3"
)

var (
	ErrInvalidPathSyntax  = errors.New("invalid path syntax")
	ErrWildCardNotAllowed = errors.New("wildcard notation not allowed for build")
)

type Path struct {
	Path string
	f    buildFn
}

type Paths []*Path

type buildFn func() *yaml.Node

func newBuilder(f func() *yaml.Node) *Path {
	return &Path{f: f}
}

// NewPath constructs a Path from a string expression.
func NewPath(path string) (*Path, error) {
	builder, err := newPath(lex("Path lexer", path))
	if err != nil {
		return nil, err
	}

	builder.Path = path

	return builder, nil
}

func NewPaths(paths []string) (*Paths, error) {
	results := make(Paths, len(paths))

	for i, p := range paths {
		np, err := NewPath(p)
		if err != nil {
			return nil, err
		}

		results[i] = np
	}

	return &results, nil
}

func newPath(l *lexer) (*Path, error) {
	lx := l.nextLexeme()

	switch lx.typ {
	case lexemeError:
		return nil, fmt.Errorf("%w, %s", ErrInvalidPathSyntax, lx.val)

	case lexemeIdentity, lexemeEOF:
		return newBuilder(identity), nil

	case lexemeRoot:
		subPath, err := newPath(l)
		if err != nil {
			return nil, err
		}

		return newBuilder(func() *yaml.Node {
			return &yaml.Node{
				Kind:        yaml.DocumentNode,
				Style:       0,
				Tag:         "",
				Value:       "",
				Anchor:      "",
				Alias:       nil,
				Content:     []*yaml.Node{subPath.f()},
				HeadComment: "",
				LineComment: "",
				FootComment: "",
				Line:        0,
				Column:      0,
			}
		}), nil

	case lexemeDotChild:
		childName := strings.TrimPrefix(lx.val, ".")

		if childName == "*" {
			return nil, fmt.Errorf("%w, wildcard notation not allowed for build", ErrInvalidPathSyntax)
		}

		childName = unescape(childName)

		return childThen(l, childName)

	case lexemeUndottedChild:
		return childThen(l, lx.val)
	case lexemeBracketChild:
		childNames := strings.TrimSpace(lx.val)
		childNames = strings.TrimSuffix(strings.TrimPrefix(childNames, "["), "]")
		childNames = strings.TrimSpace(childNames)

		return bracketChildThen(l, childNames)
	}

	return nil, ErrInvalidPathSyntax
}

func (paths *Paths) BuildPaths() (*yaml.Node, error) {
	yamlNodes := make([]*yaml.Node, len(*paths))

	for i, p := range *paths {
		yamlNode, err := p.BuildPath()
		if err != nil {
			return nil, err
		}

		yamlNodes[i] = yamlNode
	}

	if err := actions.MergeNode(yamlNodes...); err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return yamlNodes[0], nil
}

// BuildPath constructs a Path from a string expression.
func (p *Path) BuildPath() (*yaml.Node, error) {
	return p.f(), nil
}

func identity() *yaml.Node {
	return &yaml.Node{
		Kind:        yaml.ScalarNode,
		Style:       0,
		Tag:         "",
		Value:       "",
		Anchor:      "",
		Alias:       nil,
		Content:     nil,
		HeadComment: "",
		LineComment: "",
		FootComment: "",
		Line:        0,
		Column:      0,
	}
}

func childThen(l *lexer, childName string) (*Path, error) {
	if childName == "*" {
		return nil, ErrWildCardNotAllowed
	}

	childName = unescape(childName)

	subPath, err := newPath(l)
	if err != nil {
		return nil, err
	}

	return newBuilder(func() *yaml.Node {
		return &yaml.Node{
			Kind:   yaml.MappingNode,
			Style:  0,
			Tag:    "",
			Value:  "",
			Anchor: "",
			Alias:  nil,
			Content: []*yaml.Node{
				{
					Kind:        yaml.ScalarNode,
					Style:       0,
					Tag:         "",
					Value:       childName,
					Anchor:      "",
					Alias:       nil,
					Content:     nil,
					HeadComment: "",
					LineComment: "",
					FootComment: "",
					Line:        0,
					Column:      0,
				},
				subPath.f(),
			},
			HeadComment: "",
			LineComment: "",
			FootComment: "",
			Line:        0,
			Column:      0,
		}
	}), nil
}

func bracketChildNames(childNames string) []string {
	s := strings.Split(childNames, ",")
	// reconstitute child names with embedded commas
	children := []string{}
	name := ""

	for _, c := range s {
		switch {
		case balanced(c, '\'') && balanced(c, '"'):
			if name != "" {
				name += "," + c
			} else {
				children = append(children, c)
				name = ""
			}
		case name == "":
			name = c
		default:
			name += "," + c
			children = append(children, name)
			name = ""
		}
	}

	if name != "" {
		children = append(children, name)
	}

	unquotedChildren := []string{}

	for _, c := range children {
		c = strings.TrimSpace(c)
		if strings.HasPrefix(c, "'") {
			c = strings.TrimSuffix(strings.TrimPrefix(c, "'"), "'")
		} else {
			c = strings.TrimSuffix(strings.TrimPrefix(c, `"`), `"`)
		}

		c = unescape(c)
		unquotedChildren = append(unquotedChildren, c)
	}

	return unquotedChildren
}

func balanced(c string, q rune) bool {
	bal := true
	prev := eof

	for i := 0; i < len(c); {
		r, width := utf8.DecodeRuneInString(c[i:])
		i += width

		if r == q {
			if i > 0 && prev == '\\' {
				prev = r

				continue
			}

			bal = !bal
		}

		prev = r
	}

	return bal
}

func bracketChildThen(l *lexer, childNames string) (*Path, error) {
	unquotedChildren := bracketChildNames(childNames)

	subPath, err := newPath(l)
	if err != nil {
		return nil, err
	}

	return newBuilder(func() *yaml.Node {
		node := &yaml.Node{
			Kind:        yaml.MappingNode,
			Style:       0,
			Tag:         "",
			Value:       "",
			Anchor:      "",
			Alias:       nil,
			Content:     []*yaml.Node{},
			HeadComment: "",
			LineComment: "",
			FootComment: "",
			Line:        0,
			Column:      0,
		}

		for _, childName := range unquotedChildren {
			content := []*yaml.Node{
				{
					Kind:        yaml.ScalarNode,
					Style:       0,
					Tag:         "",
					Value:       childName,
					Anchor:      "",
					Alias:       nil,
					Content:     nil,
					HeadComment: "",
					LineComment: "",
					FootComment: "",
					Line:        0,
					Column:      0,
				},
				subPath.f(),
			}

			node.Content = append(node.Content, content...)
		}

		return node
	}), nil
}

func unescape(raw string) string {
	esc := ""
	escaped := false

	for i := 0; i < len(raw); {
		r, width := utf8.DecodeRuneInString(raw[i:])
		i += width

		if r == '\\' {
			if escaped {
				esc += string(r)
			}

			escaped = !escaped

			continue
		}

		escaped = false
		esc += string(r)
	}

	return esc
}
