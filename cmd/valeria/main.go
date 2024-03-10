package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/Gornak40/algolymp/config"
	"github.com/Gornak40/algolymp/polygon"
	"github.com/Gornak40/algolymp/polygon/valeria"
	"github.com/Gornak40/algolymp/polygon/valeria/textables"
	"github.com/akamensky/argparse"
	"github.com/sirupsen/logrus"
)

var (
	ErrUnknownTexTable = errors.New("unknown textable")
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
	tableTyp := parser.Selector("t", "table-type", []string{
		textables.UniversalTag,
	}, &argparse.Options{
		Required: false,
		Default:  textables.UniversalTag,
		Help:     "Tex table type",
	})
	cntVars := parser.Int("c", "count-depvars", &argparse.Options{
		Required: false,
		Default:  0,
		Help:     "Depvars count (useful for some textables)",
	})
	if err := parser.Parse(os.Args); err != nil {
		logrus.WithError(err).Fatal("bad arguments")
	}

	cfg := config.NewConfig()
	pClient := polygon.NewPolygon(&cfg.Polygon)
	val := valeria.NewValeria(pClient)

	table := textables.GetTexTable(*tableTyp, *cntVars)
	if table == nil {
		logrus.WithError(ErrUnknownTexTable).Fatal("failed to get textable")
	}

	if err := val.InformaticsValuer(*pID, table, *verbose); err != nil {
		logrus.WithError(err).Fatal("failed to get scoring")
	}

	fmt.Println(table.String()) //nolint:forbidigo // Basic functionality.
}
