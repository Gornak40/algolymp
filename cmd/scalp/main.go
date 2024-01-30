package main

import (
	"os"

	"github.com/Gornak40/algolymp/config"
	"github.com/Gornak40/algolymp/polygon"
	"github.com/akamensky/argparse"
	"github.com/sirupsen/logrus"
)

func main() {
	parser := argparse.NewParser("scalp", "Set incremental problem scoring using Polygon API.")
	pID := parser.Int("i", "problem_id", &argparse.Options{
		Required: true,
		Help:     "Polygon problem ID",
	})
	samples := parser.Flag("s", "samples", &argparse.Options{
		Required: false,
		Help:     "Include samples in scoring",
	})
	if err := parser.Parse(os.Args); err != nil {
		logrus.WithError(err).Fatal("bad arguments")
	}

	cfg := config.NewConfig()
	pClient := polygon.NewPolygon(&cfg.Polygon)

	if err := pClient.IncrementalScoring(*pID, *samples); err != nil {
		logrus.WithError(err).Fatal("failed set scoring")
	}
}
