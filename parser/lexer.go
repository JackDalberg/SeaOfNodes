package parser

import (
	"errors"
	"unicode"
)

type lexer struct {
	input    []byte
	position int
}

func NewLexer(source string) *lexer {
	return &lexer{
		input: []byte(source),
	}
}

func (l *lexer) isEOF() bool {
	return l.position >= len(l.input)
}

func (l *lexer) peek() (byte, bool) {
	if l.isEOF() {
		return 0, false
	}
	return l.input[l.position], true
}

func (l *lexer) next() (byte, bool) {
	c, ok := l.peek()
	l.position++
	return c, ok
}

func (l *lexer) isWhitespace() bool {
	c, ok := l.peek()
	if !ok {
		return false
	}
	return c <= ' '
}

func (l *lexer) skipWhitespace() {
	for l.isWhitespace() {
		l.position++
	}
}

func isValidIDStart(b byte) bool {
	return unicode.IsLetter(rune(b)) || b == '_'
}

func isValidIDByte(b byte) bool {
	return unicode.In(rune(b), unicode.Letter, unicode.Digit) || b == '_'
}

func (l *lexer) scanToken(isInToken func(b byte) bool) (string, int) {
	start := l.position
	for {
		b, ok := l.next()
		if !ok || !isInToken(b) {
			break
		}
	}
	l.position--
	s := string(l.input[start:l.position])
	return s, start
}

func (l *lexer) scanID() (string, int) {
	return l.scanToken(isValidIDByte)
}

func (l *lexer) scanPunctuation() (string, int) {
	return l.scanToken(func(b byte) bool {
		return unicode.IsPunct(rune(b))
	})
}

func (l *lexer) scanNumberString() (string, int, error) {
	s, pos := l.scanToken(func(b byte) bool {
		return unicode.IsDigit(rune(b))
	})

	if len(s) > 1 && s[0] == '0' {
		return "", pos, errors.New("integer value cannot start with '0'")
	}

	if len(s) == 0 {
		return "", pos, errors.New("not a number")
	}
	return s, pos, nil
}

func (l *lexer) ReadNumber() (string, int, error) {
	l.skipWhitespace()
	return l.scanNumberString()
}

func (l *lexer) ReadToken() (string, int, error) {
	l.skipWhitespace()
	b, ok := l.peek()

	switch {
	case !ok:
		return "", 0, nil

	case isValidIDStart(b):
		id, pos := l.scanID()
		return id, pos, nil

	case unicode.IsNumber(rune(b)):
		return l.scanNumberString()

	case unicode.IsPunct(rune(b)):
		p, pos := l.scanPunctuation()
		return p, pos, nil

	default:
		l.position++
		return string(b), l.position - 1, nil
	}
}

func (l *lexer) ReadByte() (byte, int, bool) {
	l.skipWhitespace()
	b, ok := l.next()
	if !ok {
		return 0, 0, false
	}
	return b, l.position - 1, true
}
