package main

import (
	"os"
	"path/filepath"

	"github.com/Gornak40/algolymp/config"
	"github.com/Gornak40/algolymp/polygon"
	"github.com/Gornak40/algolymp/polygon/wooda"
	"github.com/akamensky/argparse"
	"github.com/facette/natsort"
	"github.com/sirupsen/logrus"
)

func main() {
	woodaModes := []string{
		wooda.ModeTest,
		wooda.ModeTags,
		wooda.ModeValidator,
		wooda.ModeChecker,
		wooda.ModeInteractor,
		wooda.ModeSolutionMain,
		wooda.ModeSolutionCorrect,
		wooda.ModeSolutionIncorrect,
		wooda.ModeSample,
	}

	parser := argparse.NewParser("wooda", "Upload problem files filtered by glob to Polygon.")
	pID := parser.Int("i", "pid", &argparse.Options{
		Required: true,
		Help:     "Polygon problem ID",
	})
	mode := parser.Selector("m", "mode", woodaModes, &argparse.Options{
		Required: true,
		Help:     "Uploading mode",
	})
	glob := parser.String("g", "glob", &argparse.Options{
		Required: true,
		Help:     "Problem files glob",
	})
	if err := parser.Parse(os.Args); err != nil {
		logrus.WithError(err).Fatal("bad arguments")
	}

	cfg := config.NewConfig()
	pClient := polygon.NewPolygon(&cfg.Polygon)
	wooda := wooda.NewWooda(pClient, *pID, *mode)

	files, err := filepath.Glob(*glob)
	if err != nil {
		logrus.WithError(err).Fatal("failed to match glob")
	}
	if len(files) == 0 {
		logrus.WithField("glob", *glob).Warn("no files matched glob")

		return
	}
	natsort.Sort(files)
	logrus.WithFields(logrus.Fields{"glob": *glob, "count": len(files)}).
		Info("glob match result")

	errCount := 0
	for _, path := range files {
		if err := wooda.Resolve(path); err != nil {
			errCount++
			logrus.WithError(err).WithField("path", path).Error("failed to resolve")
		}
	}

	if errCount == 0 {
		logrus.Info("success resolve all files")
	} else {
		logrus.WithField("count", errCount).Warn("some errors happened")
	}
}
