package miniparse

import (
	"bufio"
	"errors"
	"io"
	"reflect"
)

var (
	ErrLeadingSpace    = errors.New("leading spaces are not allowed")
	ErrInvalidChar     = errors.New("invalid leading char")
	ErrExpectedNewLine = errors.New("expected new line")
	ErrInvalidSection  = errors.New("invalid section name")
	ErrRootSection     = errors.New("expected root section")
	ErrInvalidKey      = errors.New("invalid key name")
	ErrExpectedEqual   = errors.New("expected equal sign")
	ErrExpectedSpace   = errors.New("expected space")
	ErrUnexpectedEOF   = errors.New("unexpected end of file")

	ErrExpectedPointer = errors.New("expected not nil pointer")
	ErrExpectedStruct  = errors.New("expected struct")
	ErrBadRecordType   = errors.New("bad record type")
	ErrExpectedArray   = errors.New("expected array")
)

// Decode .mini config file into value using reflect.
// The mini format is similar to ini, but very strict.
func Decode(r io.Reader, v any) error {
	rb := bufio.NewReader(r)
	m := newMachine()
	nxt := m.stateInit

	for {
		c, _, err := rb.ReadRune()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return err
		}
		nxt, err = nxt(c)
		if err != nil {
			return err
		}
	}
	// TODO: find more Go-like solution
	if reflect.ValueOf(nxt).Pointer() != reflect.ValueOf(m.stateInit).Pointer() {
		return ErrUnexpectedEOF
	}

	return m.feed(v)
}
