package main

import (
	"fmt"
	"github.com/Gornak40/algolymp/config"
	"github.com/Gornak40/algolymp/ejudge"
	"github.com/akamensky/argparse"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"strconv"
)

const (
	DefaultRunsCount = 20
)

func main() {
	parser := argparse.NewParser("boban", "Filter Ejudge runs.")
	cID := parser.Int("i", "cid", &argparse.Options{
		Required: true,
		Help:     "Ejudge contest ID",
	})
	filter := parser.String("f", "filter", &argparse.Options{
		Required: false,
		Help:     "Filter expression",
	})
	count := parser.Int("c", "count", &argparse.Options{
		Required: false,
		Help:     "Last runs count",
		Default:  DefaultRunsCount,
	})
	inputSourcesDst := parser.String("d", "destination", &argparse.Options{
		Required: false,
		Help:     "Download run's source code to directory",
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

	runs, err := ejClient.FilterRuns(csid, *filter, *count)
	if err != nil {
		logrus.WithError(err).Fatal("filter runs failed")
	}

	sourcesDst := ""
	if len(*inputSourcesDst) > 0 {
		sourcesDst = makeContestDir(*inputSourcesDst, *cID)
	}

	for _, run := range runs {
		fmt.Println(run) //nolint:forbidigo // Basic functionality.
		if len(sourcesDst) > 0 {
			_, _ = downloadSourceCode(ejClient, csid, run, sourcesDst)
		}
	}

	if err := ejClient.Logout(sid); err != nil {
		logrus.WithError(err).Fatal("logout failed")
	}
}

func downloadSourceCode(ejClient *ejudge.Ejudge, csid string, runID int, dst string) (string, error) {
	filename, err := ejClient.DownloadRunFile(csid, runID, dst)
	if err != nil {
		logrus.WithError(err).Fatal("failed download run file")
	}

	return filename, err
}

func makeContestDir(dst string, cID int) string {
	dst = filepath.Join(filepath.Dir(dst), strconv.Itoa(cID))
	if _, err := os.Stat(dst); os.IsNotExist(err) {
		err := os.MkdirAll(dst, 0644) //nolint:mnd
		if err != nil {
			logrus.WithError(err).Fatal("directory create failed")
		}
		logrus.Info("directory [" + dst + "] created successfully")
	} else {
		logrus.Info("directory [" + dst + "] already exists")
	}

	return dst
}
