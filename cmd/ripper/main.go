package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/Gornak40/algolymp/config"
	"github.com/Gornak40/algolymp/ejudge"
	"github.com/akamensky/argparse"
	"github.com/sirupsen/logrus"
)

func main() {
	verdicts := make([]string, 0, len(ejudge.Verdicts))
	for v := range ejudge.Verdicts {
		verdicts = append(verdicts, v)
	}
	sort.Strings(verdicts)

	parser := argparse.NewParser("ripper", "Change Ejudge runs status (stdin input).")
	cID := parser.Int("i", "cid", &argparse.Options{
		Required: true,
		Help:     "Ejudge contest ID",
	})
	status := parser.Selector("s", "status", verdicts, &argparse.Options{
		Required: true,
		Help:     "New runs status",
	})
	comment := parser.String("c", "comment", &argparse.Options{
		Required: false,
		Help:     "Send run comment",
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

	logrus.Info("waiting for run ids input...")
	for {
		var runID int
		_, err := fmt.Scanf("%d", &runID)
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			logrus.WithError(err).Fatal("invalid run id")
		}
		if err := ejClient.ChangeRunStatus(csid, runID, *status); err != nil {
			logrus.WithError(err).Fatal("failed change run status")
		}
		if *comment != "" {
			if err := ejClient.SendRunComment(csid, runID, *comment); err != nil {
				logrus.WithError(err).Fatal("failed to send run comment")
			}
		}
	}

	if err := ejClient.Logout(sid); err != nil {
		logrus.WithError(err).Fatal("logout failed")
	}
}
