package postyk

import (
	"os"
	"path"

	"github.com/Gornak40/algolymp/internal/printer"
	"github.com/sirupsen/logrus"
)

const (
	cachePerm     = 0744
	cacheFilePerm = 0644
)

func (i *Indexer) Sync() error {
	if err := os.MkdirAll(i.cachePath, cachePerm); err != nil {
		return err
	}
	synced, err := os.ReadDir(i.cachePath)
	if err != nil {
		return err
	}
	logrus.WithFields(logrus.Fields{"count": len(synced), "path": i.cachePath}).
		Info("read cache directory")
	synMapa := make(map[string]struct{}, len(synced))
	for _, s := range synced {
		if !s.IsDir() {
			synMapa[s.Name()] = struct{}{}
		}
	}

	shared, err := i.GetList()
	if err != nil {
		return err
	}
	for _, s := range shared {
		if _, ok := synMapa[s.raw]; !ok {
			data, err := i.GetFile(s.raw)
			if err != nil {
				return err
			}
			fname := path.Join(i.cachePath, s.raw)
			if err := os.WriteFile(fname, data, cacheFilePerm); err != nil {
				return err
			}
			if i.printer == "" {
				logrus.Warn("no device provided, skip printing")

				continue
			}
			if err := printer.PrintFile(fname, i.printer); err != nil {
				return err
			}
		}
	}

	return nil
}
