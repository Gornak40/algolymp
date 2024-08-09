package vydra

import (
	"encoding/xml"
	"errors"
	"os"
	"strings"

	"github.com/Gornak40/algolymp/internal/natstream"
	"github.com/Gornak40/algolymp/polygon"
	"github.com/sirupsen/logrus"
)

const (
	megabyte = 1024 * 1024

	defaultTL     = 1000
	defaultML     = 256
	defaultInput  = "stdin"
	defaultOutput = "stdout"

	chkTests = "files/tests/checker-tests"
)

var (
	ErrBadSolutionTag = errors.New("bad solution tag")
)

type Vydra struct {
	client    *polygon.Polygon
	pID       int
	prob      ProblemXML
	streamIn  *natstream.NatStream
	streamOut *natstream.NatStream
	streamAns *natstream.NatStream
}

func NewVydra(client *polygon.Polygon, pID int) *Vydra {
	return &Vydra{
		client:    client,
		pID:       pID,
		streamIn:  new(natstream.NatStream),
		streamOut: new(natstream.NatStream),
		streamAns: new(natstream.NatStream),
	}
}

// Convert string from .xml to API.
func convertString(s string) string {
	return strings.ToUpper(strings.ReplaceAll(s, "-", "_")) // oh my God
}

func (v *Vydra) readXML(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	if err := xml.Unmarshal(data, &v.prob); err != nil {
		return err
	}
	logrus.WithFields(logrus.Fields{
		"revision": v.prob.Revision, "short-name": v.prob.ShortName,
	}).Info("load problem.xml")

	return nil
}

func (v *Vydra) Upload(errs chan error) error {
	defer close(errs)
	if err := v.readXML("problem.xml"); err != nil {
		return err
	}
	v.batchInitial(errs)
	v.batchValChk(errs)
	v.batchJudging(errs)

	return nil
}
