package main

import (
	"os"

	"github.com/Gornak40/algolymp/internal/korob"
	"github.com/akamensky/argparse"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

func main() {
	parser := argparse.NewParser("korob", "Generate PNG grid for problem statements.")
	input := parser.String("i", "input", &argparse.Options{
		Required: true,
		Help:     "Input grid YAML config",
	})
	output := parser.String("o", "output", &argparse.Options{
		Required: true,
		Help:     "PNG output path",
	})
	if err := parser.Parse(os.Args); err != nil {
		logrus.WithError(err).Fatal("bad arguments")
	}

	data, err := os.ReadFile(*input)
	if err != nil {
		logrus.WithError(err).Fatal("invalid grid config path")
	}
	var krb korob.Korob
	if err := yaml.Unmarshal(data, &krb); err != nil {
		logrus.WithError(err).Fatal("failed to unmarshal grid config")
	}
	logrus.WithField("input", *input).Info("config is loaded")
	if err := krb.Draw(*output); err != nil {
		logrus.WithError(err).Fatal("failed to draw grid")
	}
	logrus.WithField("output", *output).Info("success drawing grid")
}
