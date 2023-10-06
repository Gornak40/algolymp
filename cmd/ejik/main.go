package main

import (
	"flag"
	"strconv"

	"github.com/Gornak40/algolymp/config"
	"github.com/Gornak40/algolymp/ejudge"
	"github.com/sirupsen/logrus"
)

func main() {
	verboseFlag := flag.Bool("v", false, "show full output of check contest settings")
	flag.Parse()
	if flag.NArg() == 0 {
		logrus.Fatal("position argument cid required")
	}
	cid, err := strconv.Atoi(flag.Arg(0))
	if err != nil {
		logrus.WithField("cid", cid).Fatal("cid should be int")
	}

	cfg := config.NewConfig("config/config.json")
	ejClient := ejudge.NewEjudge(&cfg.Ejudge)

	sid, err := ejClient.Login()
	if err != nil {
		logrus.WithError(err).Fatal("login failed")
	}

	if err := ejClient.CheckContest(sid, cid, *verboseFlag); err != nil {
		logrus.WithError(err).Fatal("check failed")
	}

	if err := ejClient.Logout(sid); err != nil {
		logrus.WithError(err).Fatal("logout failed")
	}
}
