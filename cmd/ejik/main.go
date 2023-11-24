package main

import (
	"os"

	"github.com/Gornak40/algolymp/config"
	"github.com/Gornak40/algolymp/ejudge"
	"github.com/akamensky/argparse"
	"github.com/sirupsen/logrus"
)

func main() {
	parser := argparse.NewParser("ejik", "Refresh Ejudge contest by id.")
	cIDArg := parser.Int("i", "cid", &argparse.Options{
		Required: true,
		Help:     "Ejudge contest ID",
	})
	verboseArg := parser.Flag("v", "verbose", &argparse.Options{
		Required: false,
		Help:     "Show full output of check contest settings",
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

	if err := ejClient.Commit(sid, *cIDArg); err != nil {
		logrus.WithError(err).Fatal("commit failed")
	}

	if err := ejClient.CheckContest(sid, *cIDArg, *verboseArg); err != nil {
		logrus.WithError(err).Fatal("check failed")
	}

	if err := ejClient.ReloadConfig(sid, *cIDArg); err != nil {
		logrus.WithError(err).Fatal("reload config failed")
	}

	if err := ejClient.Logout(sid); err != nil {
		logrus.WithError(err).Fatal("logout failed")
	}
}
