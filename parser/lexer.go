package parser

import (
	"errors"
	"strings"
	"unicode"
)

var (
	ErrNAN = errors.New("not a number")
	ErrEOF = errors.New("EOF")
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
		return "", pos, ErrNAN
	}
	return s, pos, nil
}

func (l *lexer) ReadNumber() (string, int, error) {
	l.skipWhitespace()
	return l.scanNumberString()
}

func (l *lexer) ReadToken() (string, int, bool, error) {
	l.skipWhitespace()
	b, ok := l.peek()

	switch {
	case !ok:
		return "", l.position, false, ErrEOF

	case isValidIDStart(b):
		id, pos := l.scanID()
		return id, pos, true, nil

	case unicode.IsNumber(rune(b)):
		num, pos, err := l.scanNumberString()
		return num, pos, false, err

	default:
		l.position++
		return string(b), l.position - 1, false, nil
	}
}

func (l *lexer) ReadOp() (byte, int, bool) {
	l.skipWhitespace()
	b, ok := l.peek()
	if !ok || !strings.ContainsRune("+-*/", rune(b)) {
		return 0, 0, false
	}

	l.position++
	return b, l.position - 1, true
}

func (l *lexer) ReadByte() (byte, int, bool) {
	l.skipWhitespace()
	b, ok := l.next()
	if !ok {
		return 0, 0, false
	}
	return b, l.position - 1, true
}

func (l *lexer) ReadID() (string, int, bool) {
	l.skipWhitespace()
	b, ok := l.peek()
	if !ok || !isValidIDStart(b) {
		return "", l.position, false
	}
	id, pos := l.scanID()
	return id, pos, true
}

func (l *lexer) Read(next byte) (int, bool) {
	l.skipWhitespace()
	b, ok := l.peek()
	if !ok || b != next {
		return l.position, false
	}
	l.position++
	return l.position - 1, true
}
