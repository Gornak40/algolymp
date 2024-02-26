package main

import (
	"os"
	"path/filepath"

	"github.com/Gornak40/algolymp/config"
	"github.com/Gornak40/algolymp/polygon"
	"github.com/Gornak40/algolymp/polygon/wooda"
	"github.com/akamensky/argparse"
	"github.com/sirupsen/logrus"
)

func main() {
	cfg := config.NewConfig()
	woodaKeys := make([]string, 0, len(cfg.Polygon.Wooda))
	for k := range cfg.Polygon.Wooda {
		woodaKeys = append(woodaKeys, k)
	}

	parser := argparse.NewParser("wooda", "Upload problem files filtered by regexp to Polygon.")
	pID := parser.Int("i", "pid", &argparse.Options{
		Required: true,
		Help:     "Polygon problem ID",
	})
	mode := parser.Selector("m", "mode", woodaKeys, &argparse.Options{
		Required: true,
		Help:     "Local storage mode",
	})
	pDir := parser.String("d", "directory", &argparse.Options{
		Required: true,
		Help:     "Local storage directory",
	})
	if err := parser.Parse(os.Args); err != nil {
		logrus.WithError(err).Fatal("bad arguments")
	}

	pClient := polygon.NewPolygon(&cfg.Polygon)
	woodaCfg := cfg.Polygon.Wooda[*mode] // mode is good argparse.Selector
	wooda := wooda.NewWooda(pClient, *pID, &woodaCfg)
	if err := filepath.Walk(*pDir, wooda.DirWalker); err != nil {
		logrus.WithError(err).Fatal("failed wooda matching")
	}
}
