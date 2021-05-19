// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: MIT

package path

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

func BuildPaths(paths []string) (*yaml.Node, error) {
	yamlNodes := make([]*yaml.Node, len(paths))

	for i, path := range paths {
		yamlNode, err := BuildPath(path)
		if err != nil {
			return nil, err
		}

		yamlNodes[i] = yamlNode
	}

	if err := actions.Merge(yamlNodes...); err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return yamlNodes[0], nil
}

// BuildPath constructs a Path from a string expression.
func BuildPath(path string) (*yaml.Node, error) {
	return build(lex("Path lexer", path), nil)
}

func build(l *lexer, y *yaml.Node) (*yaml.Node, error) {
	lx := l.nextLexeme()

	switch lx.typ {
	case lexemeError:
		return nil, fmt.Errorf("%w, %s", ErrInvalidPathSyntax, lx.val)

	case lexemeIdentity, lexemeEOF:
		return identity()

	case lexemeRoot:
		if y == nil {
			subPath, err := build(l, y)
			if err != nil {
				return nil, err
			}

			return &yaml.Node{
				Kind:        yaml.DocumentNode,
				Style:       0,
				Tag:         "",
				Value:       "",
				Anchor:      "",
				Alias:       nil,
				Content:     []*yaml.Node{subPath},
				HeadComment: "",
				LineComment: "",
				FootComment: "",
				Line:        0,
				Column:      0,
			}, nil
		}

		return nil, nil

	case lexemeDotChild:
		childName := strings.TrimPrefix(lx.val, ".")

		if childName == "*" {
			return nil, fmt.Errorf("%w, wildcard notation not allowed for build", ErrInvalidPathSyntax)
		}

		childName = unescape(childName)

		return childThen(l, childName, y)

	case lexemeUndottedChild:
		return childThen(l, lx.val, y)
	case lexemeBracketChild:
		childNames := strings.TrimSpace(lx.val)
		childNames = strings.TrimSuffix(strings.TrimPrefix(childNames, "["), "]")
		childNames = strings.TrimSpace(childNames)

		return bracketChildThen(l, childNames, y)
	}

	return nil, ErrInvalidPathSyntax
}

func identity() (*yaml.Node, error) {
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
	}, nil
}

func childThen(l *lexer, childName string, y *yaml.Node) (*yaml.Node, error) {
	if childName == "*" {
		return nil, ErrWildCardNotAllowed
	}

	childName = unescape(childName)

	subPath, err := build(l, y)
	if err != nil {
		return nil, err
	}

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
			subPath,
		},
		HeadComment: "",
		LineComment: "",
		FootComment: "",
		Line:        0,
		Column:      0,
	}, nil
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

func bracketChildThen(l *lexer, childNames string, y *yaml.Node) (*yaml.Node, error) {
	unquotedChildren := bracketChildNames(childNames)

	for _, childName := range unquotedChildren {
		return childThen(l, childName, y)
	}

	return nil, nil
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
