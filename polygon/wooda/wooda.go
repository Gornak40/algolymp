package wooda

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"

	"github.com/Gornak40/algolymp/polygon"
	"github.com/sirupsen/logrus"
)

type Wooda struct {
	client    *polygon.Polygon
	pID       int
	config    *polygon.WoodaConfig
	testIndex int
}

func NewWooda(pClient *polygon.Polygon, pID int, wCfg *polygon.WoodaConfig) *Wooda {
	return &Wooda{
		client:    pClient,
		pID:       pID,
		config:    wCfg,
		testIndex: 1,
	}
}

func pathMatch(pattern, path string) bool {
	if pattern == "" {
		return false
	}
	res, err := regexp.MatchString(pattern, path)
	if err != nil {
		logrus.WithError(err).Error("failed to match filepath")

		return false
	}

	return res
}

func getData(mode, path string) (string, error) {
	logrus.WithFields(logrus.Fields{"mode": mode, "path": path}).Info("resolve file")
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func (w *Wooda) resolveTest(path string) error {
	data, err := getData("test", path)
	if err != nil {
		return err
	}

	tr := polygon.NewTestRequest(w.pID, w.testIndex).
		Input(data).
		Description(fmt.Sprintf("File \"%s\"", filepath.Base(path)))
	if err := w.client.SaveTest(tr); err != nil {
		return err
	}
	w.testIndex++

	return nil
}

func (w *Wooda) matcher(path string) error {
	switch {
	case pathMatch(w.config.Ignore, path): // silent ignore is the best practice
		break
	case pathMatch(w.config.Test, path):
		if err := w.resolveTest(path); err != nil {
			return err
		}
	default:
		logrus.WithField("path", path).Warn("no valid matching")
	}

	return nil
}

func (w *Wooda) DirWalker(path string, info fs.FileInfo, err error) error {
	if err != nil {
		return err
	}
	if info.IsDir() {
		return nil
	}
	if err := w.matcher(path); err != nil {
		logrus.WithError(err).WithField("path", path).Errorf("failed to resolve")
	}
	return nil
}
