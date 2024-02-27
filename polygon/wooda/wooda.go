package wooda

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Gornak40/algolymp/polygon"
	"github.com/sirupsen/logrus"
)

const (
	ModeTest      = "test"
	ModeTags      = "tags"
	ModeValidator = "validator"
)

var (
	ErrUnknownMode = errors.New("unknown wooda mode")
)

type Wooda struct {
	client    *polygon.Polygon
	pID       int
	mode      string
	testIndex int
}

func NewWooda(pClient *polygon.Polygon, pID int, mode string) *Wooda {
	return &Wooda{
		client:    pClient,
		pID:       pID,
		mode:      mode,
		testIndex: 1,
	}
}

func (w *Wooda) Resolve(path string) error {
	logrus.WithFields(logrus.Fields{"mode": w.mode, "path": path}).Info("resolve file")
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	switch w.mode {
	case ModeTest:
		return w.resolveTest(path, string(data))
	case ModeTags:
		return w.resolveTags(string(data))
	case ModeValidator:
		return w.resolveValidator(path, string(data))
	default:
		return fmt.Errorf("%w: %s", ErrUnknownMode, w.mode)
	}
}

func (w *Wooda) resolveTest(path, data string) error {
	tr := polygon.NewTestRequest(w.pID, w.testIndex).
		Input(data).
		Description(fmt.Sprintf("File \"%s\"", filepath.Base(path)))
	if err := w.client.SaveTest(tr); err != nil {
		return err
	}
	w.testIndex++

	return nil
}

func (w *Wooda) resolveTags(data string) error {
	tags := strings.Join(strings.Split(data, "\n"), ",")

	return w.client.SaveTags(w.pID, tags)
}

func (w *Wooda) resolveValidator(path, data string) error {
	name := filepath.Base(path)
	fr := polygon.NewFileRequest(w.pID, polygon.TypeSource, name, data)
	if err := w.client.SaveFile(fr); err != nil {
		return err
	}

	return w.client.SetValidator(w.pID, name)
}
