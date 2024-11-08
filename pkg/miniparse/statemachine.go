package miniparse

import (
	"fmt"
	"unicode"
)

const bufSize = 512

type machine struct {
	buf []rune
	sec string
	key string
	val string
}

func newMachine() *machine {
	return &machine{
		buf: make([]rune, 0, bufSize),
	}
}

type stateFunc func(c rune) (stateFunc, error)

func (m *machine) stateInit(c rune) (stateFunc, error) {
	switch {
	case c == '#':
		return m.stateComment, nil
	case c == '[':
		return m.stateSection, nil
	case isValidVar(c, true):
		m.buf = append(m.buf, c)

		return m.stateKey, nil
	case c == '\n':
		return m.stateInit, nil
	case unicode.IsSpace(c):
		return nil, ErrLeadingSpace
	default:
		return nil, fmt.Errorf("%w: %c", ErrInvalidChar, c)
	}
}

func (m *machine) stateComment(c rune) (stateFunc, error) {
	if c == '\n' {
		return m.stateInit, nil
	}

	return m.stateComment, nil
}

func (m *machine) stateSection(c rune) (stateFunc, error) {
	switch {
	case isValidVar(c, false):
		m.buf = append(m.buf, c)

		return m.stateSection, nil
	case c == ']':
		return m.stateSectionEnd, nil
	default:
		return nil, fmt.Errorf("%w, found: %c", ErrInvalidSection, c)
	}
}

func (m *machine) stateSectionEnd(c rune) (stateFunc, error) {
	if c != '\n' {
		return nil, fmt.Errorf("%w, found: %c", ErrExpectedNewLine, c)
	}
	m.sec = string(m.buf)
	println("sec:", m.sec) //nolint:forbidigo // ! remove
	m.buf = m.buf[:0]

	return m.stateInit, nil
}

func (m *machine) stateKey(c rune) (stateFunc, error) {
	switch {
	case isValidVar(c, false):
		m.buf = append(m.buf, c)

		return m.stateKey, nil
	case c == ' ':
		m.key = string(m.buf)
		println("key:", m.key) //nolint:forbidigo // ! remove
		m.buf = m.buf[:0]

		return m.stateEqualSign, nil
	default:
		return nil, fmt.Errorf("%w, found: %c", ErrInvalidKey, c)
	}
}

func (m *machine) stateEqualSign(c rune) (stateFunc, error) {
	if c != '=' {
		return nil, fmt.Errorf("%w, found: %c", ErrExpectedEqual, c)
	}

	return m.stateSpaceR, nil
}

func (m *machine) stateSpaceR(c rune) (stateFunc, error) {
	if c != ' ' {
		return nil, fmt.Errorf("%w, found: %c", ErrExpectedSpace, c)
	}

	return m.stateValue, nil
}

func (m *machine) stateValue(c rune) (stateFunc, error) {
	if c == '\n' {
		m.val = string(m.buf)
		println("val:", m.val) //nolint:forbidigo // ! remove
		m.buf = m.buf[:0]

		return m.stateInit, nil
	}
	m.buf = append(m.buf, c)

	return m.stateValue, nil
}

func isValidVar(c rune, first bool) bool {
	if c >= unicode.MaxASCII {
		return false
	}

	return c == '_' || unicode.IsLower(c) || (!first && unicode.IsDigit(c))
}
