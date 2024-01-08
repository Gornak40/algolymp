package main

import (
	"os"

	"github.com/Gornak40/algolymp/config"
	"github.com/Gornak40/algolymp/ejudge"
	"github.com/akamensky/argparse"
	"github.com/sirupsen/logrus"
)

func main() {
	parser := argparse.NewParser("casper", "Change Ejudge contest visibility.")
	cID := parser.Int("i", "cid", &argparse.Options{
		Required: true,
		Help:     "Ejudge contest ID",
	})
	visible := parser.Flag("s", "show", &argparse.Options{
		Required: false,
		Help:     "Make contest visible (invisible if flag is not set)",
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

	if *visible {
		err = ejClient.MakeVisible(sid, *cID)
	} else {
		err = ejClient.MakeInvisible(sid, *cID)
	}

	if err != nil {
		logrus.WithError(err).Fatal("change visible status failed")
	}
}
