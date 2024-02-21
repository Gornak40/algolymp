package main

import (
	"fmt"
	"os"

	"github.com/Gornak40/algolymp/config"
	"github.com/Gornak40/algolymp/ejudge"
	"github.com/akamensky/argparse"
	"github.com/sirupsen/logrus"
)

const (
	DefaultRunsCount = 20
)

func main() {
	parser := argparse.NewParser("boban", "Filter Ejudge runs.")
	cID := parser.Int("i", "cid", &argparse.Options{
		Required: true,
		Help:     "Ejudge contest ID",
	})
	filter := parser.String("f", "filter", &argparse.Options{
		Required: false,
		Help:     "Filter expression",
	})
	count := parser.Int("c", "count", &argparse.Options{
		Required: false,
		Help:     "Last runs count",
		Default:  DefaultRunsCount,
	})
	if err := parser.Parse(os.Args); err != nil {
		logrus.WithError(err).Fatal("bad arguments")
	}

	cfg := config.NewConfig()
	ejClient := ejudge.NewEjudge(&cfg.Ejudge)

	sid, err := ejClient.Login()
	if err != nil {
		logrus.WithError(err).Fatal("login failed")
	}

	csid, err := ejClient.MasterLogin(sid, *cID)
	if err != nil {
		logrus.WithError(err).Fatal("master login failed")
	}

	runs, err := ejClient.FilterRuns(csid, *filter, *count)
	if err != nil {
		logrus.WithError(err).Fatal("filter runs failed")
	}
	for _, run := range runs {
		fmt.Println(run) //nolint:forbidigo // Basic functionality.
	}
	logrus.WithField("runs", len(runs)).Info("filter result")

	if err := ejClient.Logout(sid); err != nil {
		logrus.WithError(err).Fatal("logout failed")
	}
}
