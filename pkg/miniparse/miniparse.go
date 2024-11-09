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
	ErrBadSectionType  = errors.New("bad section type")
	ErrBadRecordType   = errors.New("bad record type")
	ErrExpectedArray   = errors.New("expected array")
	ErrRequiredField   = errors.New("field marked required")
)

// Decode .mini config file into value using reflect.
// The mini format is similar to ini, but very strict.
//
// Each line of the config must be one of the following: blank, comment, record title, record field.
// Leading spaces are not allowed. End-of-line spaces are only allowed in record field lines.
// A non-empty .mini config file must end with a blank string.
// A comment must begin with a '#' character. All comments will be ignored by the parser.
// The record title must be "[title]", where title is the non-empty name of the varname.
// Varnames contain only lowercase ascii letters, digits and the '_' character.
// The first letter of the varname must not be a digit.
// The record field must have the form “key = value”, where key is a non-empty varname.
// The value contains any valid utf-8 sequence.
// Record names and keys can be non-unique. Then they will be interpreted as arrays.
//
// The mini format does not specify data types of values.
// But this decoder works only with string, int, bool and time.Duration.
// You should use `mini:"name"` tag to designate a structure field.
// You can use the `mini-required:"true"` tag for mandatory fields.
// You can use the `mini-default:"value"` tag for default values.
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
