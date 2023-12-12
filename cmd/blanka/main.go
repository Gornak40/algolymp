package main

import (
	"os"

	"github.com/Gornak40/algolymp/config"
	"github.com/Gornak40/algolymp/ejudge"
	"github.com/akamensky/argparse"
	"github.com/sirupsen/logrus"
)

func main() {
	parser := argparse.NewParser("blanka", "Create Ejudge contest from template.")
	cIDArg := parser.Int("i", "cid", &argparse.Options{
		Required: true,
		Help:     "Ejudge new contest ID",
	})
	tIDArg := parser.Int("t", "tid", &argparse.Options{
		Required: true,
		Help:     "Ejudge template ID",
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

	if err := ejClient.CreateContest(sid, *cIDArg, *tIDArg); err != nil {
		logrus.WithError(err).Fatal("create contest failed")
	}

	if err := ejClient.Commit(sid); err != nil {
		logrus.WithError(err).Fatal("commit failed")
	}

	if err := ejClient.Logout(sid); err != nil {
		logrus.WithError(err).Fatal("logout failed")
	}
}
