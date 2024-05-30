package main

import (
	"fmt"
	"os"
	"syscall"

	"github.com/Gornak40/algolymp/servecfg"
	"github.com/akamensky/argparse"
	"github.com/sirupsen/logrus"
)

func main() {
	parser := argparse.NewParser("fara", "Explorer for serve.cfg with mass modify.")
	file := parser.File("f", "file", syscall.O_RDONLY, 0666, &argparse.Options{ //nolint:mnd // 0666 is -rw-rw-rw
		Required: false,
		Default:  "/dev/stdin", // TODO: support Windows
		Help:     "Path to serve.cfg",
	})
	query := parser.StringList("q", "query", &argparse.Options{
		Required: true,
		Help:     "Field queries in fara format",
	})
	newVal := parser.String("u", "update", &argparse.Options{
		Required: false,
		Help:     "New value for existing selected fields",
	})
	delFlag := parser.Flag("d", "delete", &argparse.Options{
		Required: false,
		Help:     "Delete selected fields",
	})
	if err := parser.Parse(os.Args); err != nil {
		logrus.WithError(err).Fatal("bad arguments")
	}

	cfg := servecfg.New(file)
	matches := cfg.Query(*query...)
	logrus.WithField("count", len(matches)).Info("matched fields")

	switch {
	case *delFlag:
		cfg.Set(servecfg.Deleter, matches...)
	case *newVal != "":
		cfg.Set(*newVal, matches...)
	default:
		for _, match := range matches {
			fmt.Println(match.String()) //nolint:forbidigo // Basic functionality.
		}

		return
	}

	fmt.Print(cfg.String()) //nolint:forbidigo // Basic functionality.
}
