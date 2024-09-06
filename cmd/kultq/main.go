package main

import (
	"os"

	"github.com/Gornak40/algolymp/config"
	"github.com/Gornak40/algolymp/ejudge/kultq"
	"github.com/akamensky/argparse"
	"github.com/sirupsen/logrus"
)

func main() {
	parser := argparse.NewParser("kultq", "Anti plagiarism checker with interactive shell.")
	contestDir := parser.String("c", "contest", &argparse.Options{
		Required: true,
		Help:     "Contest directory with <problem>/<user-login>/<run>.[py|cpp]",
	})
	if err := parser.Parse(os.Args); err != nil {
		logrus.WithError(err).Fatal("bad arguments")
	}
	if err := os.Chdir(*contestDir); err != nil {
		logrus.WithError(err).Fatal("bad contest directory")
	}

	cfg := config.NewConfig()
	engine := kultq.NewEngine(&cfg.Antiplag)
	if err := engine.Run(); err != nil {
		logrus.WithError(err).Fatal("failed to run kultq")
	}
}
