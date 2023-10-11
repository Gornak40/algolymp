package main

import (
	"fmt"
	"os"

	"github.com/Gornak40/algolymp/config"
	"github.com/Gornak40/algolymp/ejudge"
	"github.com/akamensky/argparse"
	"github.com/sirupsen/logrus"
)

func main() {
	parser := argparse.NewParser("algolymp", "Algolymp contest manager")
	cIDArg := parser.Int("i", "cid", &argparse.Options{
		Required: true,
		Help:     "Ejudge contest ID",
	})
	verboseArg := parser.Flag("v", "verbose", &argparse.Options{
		Required: false,
		Help:     "Show full output of check contest settings",
	})
	confDir, _ := os.UserHomeDir()
	configArg := parser.String("c", "config", &argparse.Options{
		Required: false,
		Help:     "JSON config path",
		Default:  fmt.Sprintf("%s/.config/algolymp/config.json", confDir),
	})
	if err := parser.Parse(os.Args); err != nil {
		logrus.WithError(err).Fatal("bad arguments")
	}

	cfg := config.NewConfig(*configArg)
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
