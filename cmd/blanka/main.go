package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"

	"github.com/Gornak40/algolymp/config"
	"github.com/Gornak40/algolymp/ejudge"
	"github.com/akamensky/argparse"
	"github.com/sirupsen/logrus"
)

func editXML(cID int, cfg *config.Config) {
	logrus.WithField("CID", cID).Info("open xml config editor")
	app := cfg.System.Editor
	xmlName := fmt.Sprintf("%06d.xml", cID)
	arg0 := path.Join(cfg.Ejudge.JudgesDir, "data", "contests", xmlName)
	cmd := exec.Command(app, arg0)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		logrus.WithError(err).Fatalf("bad command: %s %s", app, arg0)
	}
}

func main() {
	parser := argparse.NewParser("blanka", "Create Ejudge contest from template.")
	cID := parser.Int("i", "cid", &argparse.Options{
		Required: true,
		Help:     "Ejudge new contest ID",
	})
	tID := parser.Int("t", "tid", &argparse.Options{
		Required: true,
		Help:     "Ejudge template ID",
	})
	editFlag := parser.Flag("e", "edit", &argparse.Options{
		Required: false,
		Help:     "Edit contest xml config",
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

	if err := ejClient.CreateContest(sid, *cID, *tID); err != nil {
		logrus.WithError(err).Fatal("create contest failed")
	}

	if err := ejClient.Commit(sid); err != nil {
		logrus.WithError(err).Fatal("commit failed")
	}

	if err := ejClient.Logout(sid); err != nil {
		logrus.WithError(err).Fatal("logout failed")
	}

	if *editFlag {
		editXML(*cID, cfg)
	}
}
