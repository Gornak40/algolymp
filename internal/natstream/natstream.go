package natstream

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/facette/natsort"
)

var (
	ErrEndStream = errors.New("no more files in natstream")
)

type NatStream struct {
	files []string
	idx   int
}

func (ns *NatStream) Init(glob string) error {
	files, err := filepath.Glob(glob)
	if err != nil {
		return err
	}
	ns.files = files
	natsort.Sort(ns.files)
	ns.idx = 0

	return nil
}

func (ns *NatStream) Next() (string, error) {
	if ns.idx == len(ns.files) {
		return "", ErrEndStream
	}
	data, err := os.ReadFile(ns.files[ns.idx])
	if err != nil {
		return "", err
	}
	ns.idx++

	return string(data), nil
}
