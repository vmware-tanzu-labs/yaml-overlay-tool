package path

import (
	"errors"
	"fmt"
	"strings"
	"unicode/utf8"

	"gopkg.in/yaml.v3"
)

var ErrInvalidPathSyntax = errors.New("invalid path syntax")

// Build constructs a Path from a string expression.
func Build(path string) (*yaml.Node, error) {
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
		return nil, fmt.Errorf("wildcard notation not allowed for build")
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
	accum := ""

	for _, c := range s {
		if balanced(c, '\'') && balanced(c, '"') {
			if accum != "" {
				accum += "," + c
			} else {
				children = append(children, c)
				accum = ""
			}
		} else {
			if accum == "" {
				accum = c
			} else {
				accum += "," + c
				children = append(children, accum)
				accum = ""
			}
		}
	}

	if accum != "" {
		children = append(children, accum)
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
