package main

import (
	"os"
	"time"

	"github.com/Gornak40/algolymp/config"
	"github.com/Gornak40/algolymp/ejudge/postyk"
	"github.com/akamensky/argparse"
	"github.com/sirupsen/logrus"
)

const (
	timeout = 5 * time.Second
)

func main() {
	parser := argparse.NewParser("postyk", "Service for printing Ejudge submits.")
	cID := parser.Int("i", "cid", &argparse.Options{
		Required: true,
		Help:     "Ejudge contest ID",
	})
	if err := parser.Parse(os.Args); err != nil {
		logrus.WithError(err).Fatal("bad arguments")
	}

	cfg := config.NewConfig()
	ind := postyk.NewIndexer(&cfg.Ejudge)

	if err := ind.Feed(*cID); err != nil {
		logrus.WithError(err).Fatal("failed to ping print shared directory")
	}
	for {
		if err := ind.Sync(); err != nil {
			logrus.WithError(err).Fatal("sync failed")
		}
		logrus.Info("success sync")
		time.Sleep(timeout)
	}
}
