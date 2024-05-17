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

func main() {
	parser := argparse.NewParser("baron", "Register users to Ejudge contest (Pending status).")
	cID := parser.Int("i", "cid", &argparse.Options{
		Required: true,
		Help:     "Ejudge contest ID",
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

	logrus.Info("waiting for logins input...")
	for {
		var login string
		_, err := fmt.Scan(&login)
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			logrus.WithError(err).Fatal("invalid login")
		}
		if err := ejClient.RegisterUser(csid, login); err != nil {
			logrus.WithError(err).Error("register user failed")
		}
	}

	if err := ejClient.Logout(sid); err != nil {
		logrus.WithError(err).Fatal("logout failed")
	}
}
