package main

import (
	"os"

	"github.com/Gornak40/algolymp/config"
	"github.com/Gornak40/algolymp/polygon"
	"github.com/Gornak40/algolymp/polygon/gibon"
	"github.com/akamensky/argparse"
	"github.com/sirupsen/logrus"
)

func main() {
	methods := []string{
		gibon.ModeCommit,
		gibon.ModeDownload,
		gibon.ModePackage,
		gibon.ModeUpdate,
	}

	parser := argparse.NewParser("gibon", "Polygon API methods multitool.")
	pID := parser.Int("i", "pid", &argparse.Options{
		Required: true,
		Help:     "Polygon problem ID",
	})
	method := parser.Selector("m", "method", methods, &argparse.Options{
		Required: true,
		Help:     "Polygon method",
	})
	if err := parser.Parse(os.Args); err != nil {
		logrus.WithError(err).Fatal("bad arguments")
	}

	cfg := config.NewConfig()
	pClient := polygon.NewPolygon(&cfg.Polygon)
	gib := gibon.NewGibon(pClient, *pID)

	if err := gib.Resolve(*method); err != nil {
		logrus.WithError(err).Fatal("failed to resolve")
	}
	logrus.WithFields(logrus.Fields{"problem": *pID, "method": *method}).Info("success")
}
