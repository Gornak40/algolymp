package main

import (
	"os"

	"github.com/Gornak40/algolymp/config"
	"github.com/Gornak40/algolymp/polygon"
	"github.com/akamensky/argparse"
	"github.com/sirupsen/logrus"
)

const (
	modeBuild     = "build"
	modeBuildFull = "build-full"
	modeCommit    = "commit"
	modeUpdate    = "update"
)

func main() {
	modes := []string{
		modeBuild,
		modeBuildFull,
		modeCommit,
		modeUpdate,
	}

	parser := argparse.NewParser("gibon", "Polygon API methods multitool.")
	pID := parser.Int("i", "pid", &argparse.Options{
		Required: true,
		Help:     "Polygon problem ID",
	})
	mode := parser.Selector("m", "method", modes, &argparse.Options{
		Required: true,
		Help:     "Polygon method",
	})
	if err := parser.Parse(os.Args); err != nil {
		logrus.WithError(err).Fatal("bad arguments")
	}

	cfg := config.NewConfig()
	pClient := polygon.NewPolygon(&cfg.Polygon)

	switch *mode {
	case modeBuild:
		if err := pClient.BuildPackage(*pID, false, true); err != nil {
			logrus.WithError(err).Fatal("failed to build package")
		}
	case modeBuildFull:
		if err := pClient.BuildPackage(*pID, true, true); err != nil {
			logrus.WithError(err).Fatal("failed to build full package")
		}
	case modeCommit:
		if err := pClient.Commit(*pID, true, ""); err != nil {
			logrus.WithError(err).Fatal("failed to commit")
		}
	case modeUpdate:
		if err := pClient.UpdateWorkingCopy(*pID); err != nil {
			logrus.WithError(err).Fatal("failed to update working copy")
		}
	}
	logrus.WithFields(logrus.Fields{"problem": *pID, "mode": *mode}).Info("success")
}
