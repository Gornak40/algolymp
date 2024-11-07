package miniparse

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"unicode"
)

const bufSize = 512

type mode int

const (
	modeInit mode = iota
	modeComment
	modeSection
	modeSectionEnd
	modeKey
	modeEqual
	modeSpaceR
	modeValue
)

var (
	ErrLeadingSpace    = errors.New("leading spaces are not allowed")
	ErrInvalidChar     = errors.New("invalid leading char")
	ErrExpectedNewLine = errors.New("expected new line")
	ErrInvalidSection  = errors.New("invalid section name")
	ErrInvalidKey      = errors.New("invalid key name")
	ErrExpectedEqual   = errors.New("expected equal sign")
	ErrExpectedSpace   = errors.New("expected space")
)

// Decode .mini config file into value using reflect.
// The mini format is similar to ini, but very strict.
func Decode(reader io.Reader, _ any) error { //nolint:funlen,gocognit,cyclop // TODO: refactor
	r := bufio.NewReader(reader)
	var m mode
	var sec, key, val string
	buf := make([]rune, 0, bufSize)

	for {
		c, _, err := r.ReadRune()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return err
		}
		switch m {
		case modeInit:
			switch {
			case c == '#':
				m = modeComment
			case c == '[':
				m = modeSection
			case isValidVar(c, true):
				m = modeKey
				buf = append(buf, c)
			case c == '\n':
			case unicode.IsSpace(c):
				return ErrLeadingSpace
			default:
				return fmt.Errorf("%w: %c", ErrInvalidChar, c)
			}
		case modeComment:
			if c == '\n' {
				m = modeInit
			}
		case modeSection:
			switch {
			case isValidVar(c, false):
				buf = append(buf, c)
			case c == ']':
				m = modeSectionEnd
			default:
				return fmt.Errorf("%w, found: %c", ErrInvalidSection, c)
			}
		case modeSectionEnd:
			if c != '\n' {
				return fmt.Errorf("%w, found: %c", ErrExpectedNewLine, c)
			}
			m = modeInit
			sec = string(buf)
			buf = buf[:0]
			println("section:", sec) //nolint:forbidigo // ! remove
		case modeKey:
			switch {
			case isValidVar(c, false):
				buf = append(buf, c)
			case c == ' ':
				m = modeEqual
				key = string(buf)
				println("key:", key) //nolint:forbidigo // ! remove
				buf = buf[:0]
			default:
				return fmt.Errorf("%w, found: %c", ErrInvalidKey, c)
			}
		case modeEqual:
			if c != '=' {
				return fmt.Errorf("%w, found: %c", ErrExpectedEqual, c)
			}
			m = modeSpaceR
		case modeSpaceR:
			if c != ' ' {
				return fmt.Errorf("%w, found: %c", ErrExpectedSpace, c)
			}
			m = modeValue
		case modeValue:
			if c == '\n' {
				m = modeInit
				val = string(buf)
				println("val:", val) //nolint:forbidigo // ! remove
				buf = buf[:0]
			} else {
				buf = append(buf, c)
			}
		}
	}

	return nil
}

func isValidVar(c rune, first bool) bool {
	if c >= unicode.MaxASCII {
		return false
	}

	return c == '_' || unicode.IsLower(c) || (!first && unicode.IsDigit(c))
}
