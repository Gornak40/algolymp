package miniparse

import (
	"fmt"
	"unicode"
)

const bufSize = 512

type record map[string][]string

type machine struct {
	buf  []rune
	data map[string][]record
	cur  record
	key  string
}

func newMachine() *machine {
	return &machine{
		buf:  make([]rune, 0, bufSize),
		data: make(map[string][]record),
	}
}

type stateFunc func(c rune) (stateFunc, error)

func (m *machine) stateInit(c rune) (stateFunc, error) {
	switch {
	case c == '#':
		return m.stateComment, nil
	case c == '[':
		return m.stateSection1, nil
	case isValidVar(c, true):
		if m.cur == nil {
			return nil, fmt.Errorf("%w, got %c", ErrRootSection, c)
		}
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

func (m *machine) stateSection1(c rune) (stateFunc, error) {
	switch {
	case isValidVar(c, true):
		m.buf = append(m.buf, c)

		return m.stateSection, nil
	case c == ']':
		return nil, fmt.Errorf("%w: empty", ErrInvalidSection)
	default:
		return nil, fmt.Errorf("%w, found: %c", ErrInvalidSection, c)
	}
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
	sec := string(m.buf)
	m.cur = make(record)
	m.data[sec] = append(m.data[sec], m.cur)
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
		val := string(m.buf)
		m.cur[m.key] = append(m.cur[m.key], val)
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
