package miniparse

import (
	"fmt"
	"reflect"
	"strconv"
)

const (
	tagName = "mini"
)

func (m *machine) feed(v any) error {
	pv := reflect.ValueOf(v)
	if pv.Kind() != reflect.Pointer || pv.IsNil() {
		return ErrExpectedPointer
	}
	vv := pv.Elem()
	vt := vv.Type()
	if vv.Kind() != reflect.Struct {
		return ErrExpectedStruct
	}
	for i := range vv.NumField() {
		tf := vt.Field(i)
		vf := vv.Field(i)
		if err := m.parseField(tf, vf); err != nil {
			return err
		}
	}

	return nil
}

func (m *machine) parseField(f reflect.StructField, v reflect.Value) error {
	name := f.Tag.Get(tagName)
	if name == "" {
		return nil
	}
	r, ok := m.data[name]
	if !ok {
		return nil
	}
	if v.Kind() != reflect.Struct { // TODO: support arrays
		return fmt.Errorf("%w: field %s", ErrExpectedStruct, name)
	}

	return writeRecord(r[0], v) // TODO: fix it
}

func writeRecord(r record, v reflect.Value) error {
	t := v.Type()
	for i := range v.NumField() {
		tf := t.Field(i)
		name := tf.Tag.Get(tagName)
		if name == "" {
			continue
		}
		a, ok := r[name]
		if !ok {
			continue
		}
		if tf.Type.Kind() != reflect.Slice && len(a) > 1 {
			return fmt.Errorf("%w: field %s", ErrExpectedArray, name)
		}
		val, err := parseValue(a, tf.Type)
		if err != nil {
			return err
		}
		v.Field(i).Set(val)
	}

	return nil
}

func parseValue(a []string, t reflect.Type) (reflect.Value, error) {
	switch t {
	case reflect.TypeOf(""):
		return reflect.ValueOf(a[0]), nil
	case reflect.TypeOf([]string{}):
		return reflect.ValueOf(a), nil
	case reflect.TypeOf(0):
		x, err := toInts(a)
		if err != nil {
			return reflect.Value{}, err
		}

		return reflect.ValueOf(x[0]), nil
	case reflect.TypeOf([]int{}):
		x, err := toInts(a)
		if err != nil {
			return reflect.Value{}, err
		}

		return reflect.ValueOf(x), nil
	case reflect.TypeOf(false):
		x, err := toBools(a)
		if err != nil {
			return reflect.Value{}, err
		}

		return reflect.ValueOf(x[0]), nil
	case reflect.TypeOf([]bool{}):
		x, err := toBools(a)
		if err != nil {
			return reflect.Value{}, err
		}

		return reflect.ValueOf(x), nil
	default:
		return reflect.Value{}, fmt.Errorf("%w: %s", ErrBadRecordType, t.String())
	}
}

func toInts(a []string) ([]int, error) {
	ra := make([]int, 0, len(a))
	for _, s := range a {
		x, err := strconv.Atoi(s)
		if err != nil {
			return nil, err
		}
		ra = append(ra, x)
	}

	return ra, nil
}

func toBools(a []string) ([]bool, error) {
	ra := make([]bool, 0, len(a))
	for _, s := range a {
		x, err := strconv.ParseBool(s)
		if err != nil {
			return nil, err
		}
		ra = append(ra, x)
	}

	return ra, nil
}
