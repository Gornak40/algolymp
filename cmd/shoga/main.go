package main

import (
	"io"
	"os"

	"github.com/Gornak40/algolymp/config"
	"github.com/Gornak40/algolymp/ejudge"
	"github.com/akamensky/argparse"
	"github.com/sirupsen/logrus"
)

const (
	modeUsers     = "usr"
	modeRuns      = "run"
	modeStandings = "stn"
	modeProblems  = "prb"
	modePasswords = "reg"
)

func main() {
	parser := argparse.NewParser("shoga", "Dump Ejudge contest tables.")
	cID := parser.Int("i", "cid", &argparse.Options{
		Required: true,
		Help:     "Ejudge contest ID",
	})
	av := []string{modeUsers, modeRuns, modeStandings, modeProblems, modePasswords}
	mode := parser.Selector("m", "mode", av, &argparse.Options{
		Required: true,
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

	var call func(csid string) (io.Reader, error)
	switch *mode {
	case modeUsers:
		call = ejClient.DumpUsers
	case modeRuns:
		call = ejClient.DumpRuns
	case modeStandings:
		call = ejClient.DumpStandings
	case modeProblems:
		call = ejClient.DumpProbStats
	case modePasswords:
		call = ejClient.DumpRegPasswords
	}
	r, err := call(csid)
	if err != nil {
		logrus.WithField("mode", *mode).Fatal("dump failed")
	}
	if _, err := io.Copy(os.Stdout, r); err != nil {
		logrus.WithError(err).Fatal("write dumped content failed")
	}

	if err := ejClient.Logout(sid); err != nil {
		logrus.WithError(err).Fatal("logout failed")
	}
}
