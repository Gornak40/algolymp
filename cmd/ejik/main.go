package main

import (
	"github.com/Gornak40/algolymp/config"
	"github.com/Gornak40/algolymp/ejudge"
	"github.com/sirupsen/logrus"
)

func main() {
	cfg := config.NewConfig("config/config.json")
	ejClient := ejudge.NewEjudge(&cfg.Ejudge)
	sid, err := ejClient.Login()
	if err != nil {
		logrus.WithError(err).Fatal("login failed")
	}
	if err := ejClient.CheckContest(sid, "40309"); err != nil {
		logrus.WithError(err).Fatal("check failed")
	}
	if err := ejClient.Logout(sid); err != nil {
		logrus.WithError(err).Fatal("logout failed")
	}
}
