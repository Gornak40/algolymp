package main

import (
	"os"

	"github.com/Gornak40/algolymp/config"
	"github.com/Gornak40/algolymp/polygon"
	"github.com/Gornak40/algolymp/polygon/vydra"
	"github.com/akamensky/argparse"
	"github.com/sirupsen/logrus"
)

func main() {
	parser := argparse.NewParser("vydra", "Upload package to Polygon.")
	pID := parser.Int("i", "pid", &argparse.Options{
		Required: true,
		Help:     "Polygon problem ID",
	})
	pDir := parser.String("p", "prob-dir", &argparse.Options{
		Required: false,
		Default:  ".",
		Help:     "Problem directory (with problem.xml)",
	})

	if err := parser.Parse(os.Args); err != nil {
		logrus.WithError(err).Fatal("bad arguments")
	}
	if err := os.Chdir(*pDir); err != nil {
		logrus.WithError(err).Fatal("bad problem directory")
	}

	cfg := config.NewConfig()
	pClient := polygon.NewPolygon(&cfg.Polygon)

	vyd := vydra.NewVydra(pClient, *pID)
	errs := make(chan error)
	go func() {
		for err := range errs {
			if err != nil {
				logrus.WithError(err).Error("vydra error")
			}
		}
	}()
	if err := vyd.Upload(errs); err != nil {
		logrus.WithError(err).Fatal("upload failed")
	}
}
