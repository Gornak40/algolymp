package main

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/Gornak40/algolymp/config"
	"github.com/Gornak40/algolymp/ejudge"
	"github.com/akamensky/argparse"
	"github.com/sirupsen/logrus"
)

const (
	makeInvisible = "hide"
	makeVisible   = "show"
)

func main() {
	parser := argparse.NewParser("casper", "Change Ejudge contests visibility (stdin input).")
	mode := parser.Selector("m", "mode", []string{makeInvisible, makeVisible}, &argparse.Options{
		Required: true,
		Help:     "Invisible or visible",
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

	var casperFunc func(string, int) error
	switch *mode {
	case makeVisible:
		casperFunc = ejClient.MakeVisible
	case makeInvisible:
		casperFunc = ejClient.MakeInvisible
	default:
		logrus.WithField("mode", *mode).Fatal("unknown mode")
	}

	logrus.Info("waiting for contest ids input...")
	for {
		var cid int
		_, err := fmt.Scan(&cid)
		if errors.Is(err, io.EOF) {
			break
		}
		if err := casperFunc(sid, cid); err != nil {
			logrus.WithError(err).Fatal("change visible status failed")
		}
	}

	if err := ejClient.Logout(sid); err != nil {
		logrus.WithError(err).Fatal("logout failed")
	}
}
