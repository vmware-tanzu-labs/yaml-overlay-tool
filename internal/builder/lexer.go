// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: MIT

package builder

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

// This lexer was based on Rob Pike's talk "Lexical Scanning in Go" (https://talks.golang.org/2011/lex.slide#1)

type lexemeType int

const (
	lexemeError lexemeType = iota
	lexemeIdentity
	lexemeRoot
	lexemeDotChild
	lexemeUndottedChild
	lexemeBracketChild
	lexemeEOF // lexing complete
)

const (
	chanLen = 2
)

// a lexeme is a token returned from the lexer.
type lexeme struct {
	typ lexemeType
	val string // original lexeme or error message if typ is lexemeError
}

// stateFn represents the state of the lexer as a function that returns the next state.
// A nil stateFn indicates lexing is complete.
type stateFn func(*lexer) stateFn

// lexer holds the state of the scanner.
type lexer struct {
	name                  string      // name of the lexer, used only for error reports
	input                 string      // the string being scanned
	start                 int         // start position of this item
	pos                   int         // current position in the input
	width                 int         // width of last rune read from input
	state                 stateFn     // lexer state
	stack                 []stateFn   // lexer stack
	items                 chan lexeme // channel of scanned lexemes
	lastEmittedStart      int         // start position of last scanned lexeme
	lastEmittedLexemeType lexemeType  // type of last emitted lexeme (or lexemEOF if no lexeme has been emitted)
}

// lex creates a new scanner for the input string.
func lex(name, input string) *lexer {
	l := &lexer{
		name:                  name,
		input:                 input,
		state:                 lexPath,
		stack:                 make([]stateFn, 0),
		items:                 make(chan lexeme, chanLen),
		lastEmittedLexemeType: lexemeEOF,
	}

	return l
}

// pop pops a state function from the stack. If the stack is empty, returns an error function.
func (l *lexer) pop() stateFn {
	if len(l.stack) == 0 {
		return l.errorf("syntax error")
	}

	index := len(l.stack) - 1
	element := l.stack[index]
	l.stack = l.stack[:index]

	return element
}

// empty returns true if and onl if the stack of state functions is empty.
func (l *lexer) emptyStack() bool {
	return len(l.stack) == 0
}

// nextLexeme returns the next item from the input.
func (l *lexer) nextLexeme() lexeme {
	for {
		select {
		case item := <-l.items:
			return item
		default:
			if l.state == nil {
				return lexeme{
					typ: lexemeEOF,
				}
			}

			l.state = l.state(l)
		}
	}
}

const eof rune = -1 // invalid Unicode code point

// next returns the next rune in the input.
func (l *lexer) next() (r rune) {
	if l.pos >= len(l.input) {
		l.width = 0

		return eof
	}

	r, l.width = utf8.DecodeRuneInString(l.input[l.pos:])
	l.pos += l.width

	return r
}

// consume consumes as many runes as there are in the given string.
func (l *lexer) consume(s string) {
	for range s {
		l.next()
	}
}

// consumed checks the input to see if it starts with the given token. If so,
// it consumes the given token and returns true. Otherwise, it returns false.
func (l *lexer) consumed(token string) bool {
	if l.hasPrefix(token) {
		l.consume(token)

		return true
	}

	return false
}

// consumedWhitespaces checks the input to see if, after whitespace is removed, it
// starts with the given tokens. If so, it consumes the given
// tokens and any whitespace and returns true. Otherwise, it returns false.
func (l *lexer) consumedWhitespaced(tokens ...string) bool {
	pos := l.pos

	for _, token := range tokens {
		// skip past whitespace
		for {
			if pos >= len(l.input) {
				return false
			}

			r, width := utf8.DecodeRuneInString(l.input[pos:])
			if !unicode.IsSpace(r) {
				break
			}

			pos += width
		}

		if !strings.HasPrefix(l.input[pos:], token) {
			return false
		}

		pos += len(token)
	}

	l.pos = pos

	return true
}

// consumeWhitespace consumes any leading whitespace.
func (l *lexer) consumeWhitespace() {
	pos := l.pos

	for {
		if pos >= len(l.input) {
			break
		}

		r, width := utf8.DecodeRuneInString(l.input[pos:])
		if !unicode.IsSpace(r) {
			break
		}

		pos += width
	}

	l.pos = pos
}

// consumeUntil consumes tokens until it hits one of the exceptions provided,
// if no exception is provided it will consume tokens until end of input.
func (l *lexer) consumeUntil(except ...rune) bool {
	consumed := false

	except = append(except, eof)

	for {
		le := l.next()
		for _, s := range except {
			if le == s {
				l.backup()

				return consumed
			}
		}

		consumed = true
	}
}

// peeked checks the input to see if it starts with the given token and does
// not start with any of the given exceptions. If so, it returns true.
// Otherwise, it returns false.
func (l *lexer) peeked(token string, except ...string) bool {
	if l.hasPrefix(token) {
		for _, e := range except {
			if l.hasPrefix(e) {
				return false
			}
		}

		return true
	}

	return false
}

// peekedWhitespaces checks the input to see if, after whitespace is removed, it
// starts with the given tokens. If so, it returns true. Otherwise, it returns false.
func (l *lexer) peekedWhitespaced(tokens ...string) bool {
	pos := l.pos

	for _, token := range tokens {
		// skip past whitespace
		for {
			if pos >= len(l.input) {
				return false
			}

			r, width := utf8.DecodeRuneInString(l.input[pos:])
			if !unicode.IsSpace(r) {
				break
			}

			pos += width
		}

		if !strings.HasPrefix(l.input[pos:], token) {
			return false
		}

		pos += len(token)
	}

	return true
}

// backup steps back one rune.
// Can be called only once per call of next.
func (l *lexer) backup() {
	l.pos -= l.width
}

// emit passes a lexeme back to the client.
func (l *lexer) emit(typ lexemeType) {
	l.items <- lexeme{
		typ: typ,
		val: l.value(),
	}

	l.lastEmittedStart = l.start
	l.start = l.pos
	l.lastEmittedLexemeType = typ
}

// value returns the portion of the current lexeme scanned so far.
func (l *lexer) value() string {
	return l.input[l.start:l.pos]
}

// context returns the last emitted lexeme (if any) followed by the portion
// of the current lexeme scanned so far.
func (l *lexer) context() string {
	return l.input[l.lastEmittedStart:l.pos]
}

// emitSynthetic passes a lexeme back to the client which wasn't encountered in the input.
// The lexing position is not modified.
func (l *lexer) emitSynthetic(typ lexemeType, val string) {
	l.items <- lexeme{
		typ: typ,
		val: val,
	}
}

func (l *lexer) empty() bool {
	return l.pos >= len(l.input)
}

func (l *lexer) hasPrefix(p string) bool {
	return strings.HasPrefix(l.input[l.pos:], p)
}

// errorf returns an error lexeme with context and terminates the scan.
func (l *lexer) errorf(format string, args ...interface{}) stateFn {
	l.items <- lexeme{
		typ: lexemeError,
		val: fmt.Sprintf("%s at position %d, following %q", fmt.Sprintf(format, args...), l.pos, l.context()),
	}

	return nil
}

const (
	root string = "$"
	dot  string = "."
)

func lexPath(l *lexer) stateFn {
	if l.empty() {
		l.emit(lexemeIdentity)
		l.emit(lexemeEOF)

		return nil
	}

	if l.hasPrefix(root) {
		return lexRoot
	}

	// emit implicit root
	l.emitSynthetic(lexemeRoot, root)

	return lexSubPath
}

func lexRoot(l *lexer) stateFn {
	l.pos += len(root)
	l.emit(lexemeRoot)

	return lexSubPath
}

// consumedEscapedString consumes a string with the given string validly escaped using "\" and returns
// true if and only if such a string was consumed.
func consumedEscapedString(l *lexer, quote string) bool {
	for {
		switch {
		case l.peeked(quote): // unescaped quote
			return true
		case l.consumed(`\` + quote):
		case l.consumed(`\\`):
		case l.peeked(`\`):
			l.errorf("unsupported escape sequence inside %s%s", quote, quote)

			return false
		default:
			if l.next() == eof {
				l.errorf("unmatched %s", enquote(quote))

				return false
			}
		}
	}
}

func lexSubPath(l *lexer) stateFn {
	switch {
	case l.hasPrefix(")"):
		return l.pop()

	case l.empty():
		return l.handleEmpty()

	case l.consumed(dot):
		return l.handleChild(lexemeDotChild)

	case l.peekedWhitespaced("[", "'") || l.peekedWhitespaced("[", `"`): // bracketQuote or bracketDoubleQuote
		return l.handleBracket()

	case l.lastEmittedLexemeType == lexemeEOF:
		return l.handleChild(lexemeUndottedChild)
	default:
		return l.errorf("invalid path syntax")
	}
}

func enquote(quote string) string {
	switch quote {
	case "'":
		return `"'"`

	case `"`:
		return `'"'`

	default:
		panic(fmt.Sprintf(`enquote called with incorrect argument %q`, quote))
	}
}

func (l *lexer) handleEmpty() stateFn {
	if !l.emptyStack() {
		return l.pop()
	}

	l.emit(lexemeIdentity)
	l.emit(lexemeEOF)

	return nil
}

func (l *lexer) handleChild(t lexemeType) stateFn {
	exceptions := []rune{
		'.', '[', ')', ' ', '&',
		'|', '=', '!', '>', '<',
		'~', '*', eof,
	}

	if childName := l.consumeUntil(exceptions...); !childName {
		return l.errorf("child name missing")
	}

	l.emit(t)

	return lexSubPath
}

func (l *lexer) handleBracket() stateFn {
	l.consumedWhitespaced("[")

	for {
		l.consumeWhitespace()
		quote := string(l.next())

		if !consumedEscapedString(l, quote) {
			return nil
		}

		if !l.consumed(quote) {
			return l.errorf(`missing %s`, enquote(quote))
		}

		if l.consumedWhitespaced(",") {
			if !l.peekedWhitespaced("'") && !l.peekedWhitespaced(`"`) {
				return l.errorf(`missing %s or %s`, enquote("'"), enquote(`"`))
			}
		} else {
			break
		}
	}

	if !l.consumedWhitespaced("]") {
		return l.errorf(`missing "]" or ","`)
	}

	l.emit(lexemeBracketChild)

	return lexSubPath
}
