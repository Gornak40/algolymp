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

const (
	universalTag = "universal"
	moscowTag    = "moscow"
)

var (
	ErrUnknownTexTable = errors.New("unknown textable")
)

func main() {
	parser := argparse.NewParser("valeria", "Build valuer + textable using Polygon API.")
	pID := parser.Int("i", "problem_id", &argparse.Options{
		Required: true,
		Help:     "Polygon problem ID",
	})
	verbose := parser.Flag("v", "verbose", &argparse.Options{
		Required: false,
		Help:     "Print valuer.cfg in stderr",
	})
	tableTyp := parser.Selector("t", "textable-type", []string{
		universalTag,
		moscowTag,
	}, &argparse.Options{
		Required: false,
		Default:  universalTag,
		Help:     "Textable type",
	})
	vars := parser.StringList("c", "variable", &argparse.Options{
		Required: false,
		Default:  nil,
		Help:     "Variables list (useful for some textables)",
	})
	if err := parser.Parse(os.Args); err != nil {
		logrus.WithError(err).Fatal("bad arguments")
	}

	cfg := config.NewConfig()
	pClient := polygon.NewPolygon(&cfg.Polygon)
	val := valeria.NewValeria(pClient)

	logrus.WithFields(logrus.Fields{"type": *tableTyp, "vars": *vars}).Info("select textable")
	var table textables.Table
	switch *tableTyp {
	case universalTag:
		table = new(textables.UniversalTable)
	case moscowTag:
		table = textables.NewMoscowTable(*vars)
	}
	if table == nil {
		logrus.WithError(ErrUnknownTexTable).Fatal("failed to get textable")
	}

	if err := val.InformaticsValuer(*pID, table, *verbose); err != nil {
		logrus.WithError(err).Fatal("failed to get scoring")
	}

	fmt.Println(table.String()) //nolint:forbidigo // Basic functionality.
}
