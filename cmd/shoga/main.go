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
	modeUsers     = "usr"
	modeRuns      = "run"
	modeStandings = "snd"
)

func main() {
	parser := argparse.NewParser("shoga", "Dump Ejudge contest users.")
	cID := parser.Int("i", "cid", &argparse.Options{
		Required: true,
		Help:     "Ejudge contest ID",
	})
	mode := parser.Selector("m", "mode", []string{modeUsers, modeRuns, modeStandings}, &argparse.Options{
		Required: false,
		Default:  modeUsers,
		Help:     "Dump mode",
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

	var call func(csid string) (string, error)
	switch *mode {
	case modeUsers:
		call = ejClient.DumpUsers
	case modeRuns:
		call = ejClient.DumpRuns
	case modeStandings:
		call = ejClient.DumpStandings
	}
	list, err := call(csid)
	if err != nil {
		logrus.WithField("mode", *mode).Fatal("dump failed")
	}
	fmt.Print(list) //nolint:forbidigo // Basic functionality.

	if err := ejClient.Logout(sid); err != nil {
		logrus.WithError(err).Fatal("logout failed")
	}
}
