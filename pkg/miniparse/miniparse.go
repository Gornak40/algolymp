package miniparse

import (
	"bufio"
	"errors"
	"io"
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
func Decode(reader io.Reader, _ any) error {
	r := bufio.NewReader(reader)
	m := newMachine()
	nxt := m.stateInit

	for {
		c, _, err := r.ReadRune()
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

	return nil
}
