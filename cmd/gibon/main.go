package main

import (
	"os"

	"github.com/Gornak40/algolymp/config"
	"github.com/Gornak40/algolymp/polygon"
	"github.com/akamensky/argparse"
	"github.com/sirupsen/logrus"
)

const (
	modeCommit = "commit"
	modeUpdate = "update"
)

func main() {
	modes := []string{
		modeCommit,
		modeUpdate,
	}

	parser := argparse.NewParser("gibon", "Polygon API multitool.")
	pID := parser.Int("i", "pid", &argparse.Options{
		Required: true,
		Help:     "Polygon problem ID",
	})
	mode := parser.Selector("m", "mode", modes, &argparse.Options{
		Required: true,
		Help:     "Polygon method",
	})
	if err := parser.Parse(os.Args); err != nil {
		logrus.WithError(err).Fatal("bad arguments")
	}

	cfg := config.NewConfig()
	pClient := polygon.NewPolygon(&cfg.Polygon)

	switch *mode {
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
