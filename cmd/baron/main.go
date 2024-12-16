package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/Gornak40/algolymp/config"
	"github.com/Gornak40/algolymp/ejudge"
	"github.com/akamensky/argparse"
	"github.com/sirupsen/logrus"
)

const (
	modeFlipVisible = "invis"
	modeFlipBan     = "ban"
	modeFlipLock    = "lock"
	modeFlipIncom   = "incom"
	modeFlipPriv    = "priv"
)

func main() {
	parser := argparse.NewParser("baron", "Ejudge contest users manager.")
	cID := parser.Int("i", "cid", &argparse.Options{
		Required: true,
		Help:     "Ejudge contest ID",
	})
	av := []string{modeFlipVisible, modeFlipBan, modeFlipLock, modeFlipIncom, modeFlipPriv}
	mode := parser.Selector("f", "flip", av, &argparse.Options{
		Required: true,
		Help:     "Users processing flip mode",
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

	logrus.Info("waiting for user ids input...")
	for {
		var suid string
		_, err := fmt.Scan(&suid)
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			logrus.WithError(err).Fatal("scan failed")
		}
		uid, err := strconv.Atoi(suid)
		if err != nil {
			logrus.WithError(err).WithField("uid", suid).Error("invalid user id")

			continue
		}
		call := getFunc(ejClient, *mode)
		if err := call(csid, uid); err != nil {
			logrus.WithError(err).Error("processing user failed")
		}
	}

	if err := ejClient.Logout(sid); err != nil {
		logrus.WithError(err).Fatal("logout failed")
	}
}

func getFunc(ej *ejudge.Ejudge, mode string) func(string, int) error {
	switch mode {
	case modeFlipVisible:
		return ej.FlipUserVisible
	case modeFlipBan:
		return ej.FlipUserBan
	case modeFlipLock:
		return ej.FlipUserLock
	case modeFlipIncom:
		return ej.FlipUserIncom
	case modeFlipPriv:
		return ej.FlipUserPriv
	default:
		logrus.WithField("mode", mode).Fatal("unknown mode")

		return nil
	}
}
