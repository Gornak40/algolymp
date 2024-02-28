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
	ModeTest              = "t"
	ModeTags              = "tags"
	ModeValidator         = "v"
	ModeChecker           = "c"
	ModeInteractor        = "i"
	ModeSolutionMain      = "ma"
	ModeSolutionCorrect   = "ok"
	ModeSolutionIncorrect = "rj"
	ModeSample            = "s"
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
	file := string(data)
	switch w.mode {
	case ModeTest:
		return w.resolveTest(path, file, false)
	case ModeTags:
		return w.resolveTags(file)
	case ModeValidator:
		return w.resolveValidator(path, file)
	case ModeChecker:
		return w.resolveChecker(path, file)
	case ModeInteractor:
		return w.resolveInteractor(path, file)
	case ModeSolutionMain:
		return w.resolveSolution(path, file, polygon.TagMain)
	case ModeSolutionCorrect:
		return w.resolveSolution(path, file, polygon.TagCorrect)
	case ModeSolutionIncorrect:
		return w.resolveSolution(path, file, polygon.TagIncorrect)
	case ModeSample:
		return w.resolveTest(path, file, true)
	default:
		return fmt.Errorf("%w: %s", ErrUnknownMode, w.mode)
	}
}

func (w *Wooda) resolveTest(path, data string, sample bool) error {
	tr := polygon.NewTestRequest(w.pID, w.testIndex).
		Input(data).
		Description(fmt.Sprintf("File \"%s\"", filepath.Base(path))).
		UseInStatements(sample)
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

// TODO: support standard checkers.
func (w *Wooda) resolveChecker(path, data string) error {
	name := filepath.Base(path)
	fr := polygon.NewFileRequest(w.pID, polygon.TypeSource, name, data)
	if err := w.client.SaveFile(fr); err != nil {
		return err
	}

	return w.client.SetChecker(w.pID, name)
}

func (w *Wooda) resolveInteractor(path, data string) error {
	pr := polygon.NewProblemRequest(w.pID).Interactive(true)
	if err := w.client.UpdateInfo(pr); err != nil {
		return err
	}

	name := filepath.Base(path)
	fr := polygon.NewFileRequest(w.pID, polygon.TypeSource, name, data)
	if err := w.client.SaveFile(fr); err != nil {
		return err
	}

	return w.client.SetInteractor(w.pID, name)
}

func (w *Wooda) resolveSolution(path, data string, tag polygon.SolutionTag) error {
	return w.client.SaveSolution(w.pID, filepath.Base(path), data, tag)
}
