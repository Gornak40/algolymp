package main

import (
	"fmt"
	"os"

	"github.com/Gornak40/algolymp/config"
	"github.com/Gornak40/algolymp/polygon"
	"github.com/Gornak40/algolymp/polygon/valeria"
	"github.com/akamensky/argparse"
	"github.com/sirupsen/logrus"
)

func main() {
	parser := argparse.NewParser("valeria", "Build valuer + scorer using Polygon API.")
	pID := parser.Int("i", "problem_id", &argparse.Options{
		Required: true,
		Help:     "Polygon problem ID",
	})
	verbose := parser.Flag("v", "verbose", &argparse.Options{
		Required: false,
		Help:     "Print valuer.cfg in stderr",
	})
	if err := parser.Parse(os.Args); err != nil {
		logrus.WithError(err).Fatal("bad arguments")
	}

	cfg := config.NewConfig()
	pClient := polygon.NewPolygon(&cfg.Polygon)
	val := valeria.NewValeria(pClient)

	table := valeria.UniversalTable{}
	if err := val.InformaticsValuer(*pID, &table, *verbose); err != nil {
		logrus.WithError(err).Fatal("failed get scoring")
	}

	fmt.Println(table.String()) //nolint:forbidigo // Basic functionality.
}
