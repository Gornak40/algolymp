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
	"strings"
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
		Help:     "Download runs source codes to directory",
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
		if sourcesDst != "" {
			filename := downloadSourceCode(ejClient, csid, run, sourcesDst)
			appendCommentsToSourceCode(ejClient, csid, run, sourcesDst, filename)
		}
	}

	if err := ejClient.Logout(sid); err != nil {
		logrus.WithError(err).Fatal("logout failed")
	}
}

func appendCommentsToSourceCode(ejClient *ejudge.Ejudge, csid string, runId int, contestDestination string, runSourceCodeFilename string) {
	currentComments, previousComments, _ := ejClient.GetAllComments(csid, runId)

	file, err := os.OpenFile(filepath.Join(contestDestination, runSourceCodeFilename), os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		logrus.WithError(err).Fatal("cannot open runId source code file")
	}
	defer file.Close()

	writeCommentSection(file, "COMMENTS FOR CURRENT RUN", currentComments)
	writeCommentSection(file, "COMMENTS FOR PREVIOUS RUN", previousComments)
}

func writeCommentSection(file *os.File, header string, comments []ejudge.Comment) {
	if comments != nil && len(comments) != 0 {
		writeStringToFilef(file, "\n\n")
		writeStringToFilef(file, "/*")
		writeStringToFilef(file, "=============== %s", header)
		for _, comment := range comments {
			writeComment(file, comment)
		}
		writeStringToFilef(file, "===============")
		writeStringToFilef(file, "*/")
	}
}

var TAB_COUNT = 0

func writeComment(file *os.File, comment ejudge.Comment) {
	writeStringToFilef(file, "[%s]: ", comment.Author)
	TAB_COUNT++
	writeStringToFilef(file, "%s", comment.Content)
	TAB_COUNT--
}

func writeStringToFilef(file *os.File, format string, args ...interface{}) {
	tabs := strings.Repeat("\t", TAB_COUNT)
	content := fmt.Sprintf(format, args...)
	content = tabs + strings.Join(strings.Split(content, "\n"), "\n"+tabs) + "\n"
	_, err := file.WriteString(content)
	if err != nil {
		logrus.WithError(err).Fatal("cannot write comment to run source code file")
	}
}

func downloadSourceCode(ejClient *ejudge.Ejudge, csid string, runID int, dst string) string {
	filename, err := ejClient.DownloadRunFile(csid, runID, dst)
	if err != nil {
		logrus.WithError(err).Fatal("failed download run file")
	}

	return filename
}

func makeContestDir(dst string, cID int) string {
	dst = filepath.Join(filepath.Dir(dst), strconv.Itoa(cID))
	if _, err := os.Stat(dst); os.IsNotExist(err) {
		err := os.MkdirAll(dst, 0775) //nolint:mnd
		if err != nil {
			logrus.WithError(err).Fatal("directory create failed")
		}
		logrus.Infof("directory [%s] created successfully", dst)
	} else {
		logrus.Infof("directory [%s] already exists", dst)
	}

	return dst
}
