package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Gornak40/algolymp/pepel"
	"github.com/akamensky/argparse"
	"github.com/sirupsen/logrus"
)

func main() {
	parser := argparse.NewParser("pepel", "Generate hasher solution based on a/ans/out files.")
	infGlob := parser.String("i", "input", &argparse.Options{
		Required: true,
		Help:     "Input files glob",
	})
	ansGlob := parser.String("a", "answer", &argparse.Options{
		Required: true,
		Help:     "Answer files glob",
	})
	zFlag := parser.Flag("z", "zip", &argparse.Options{
		Required: false,
		Help:     "Compress answers using zlib",
	})
	if err := parser.Parse(os.Args); err != nil {
		logrus.WithError(err).Fatal("bad arguments")
	}

	infAll, err := filepath.Glob(*infGlob)
	if err != nil {
		logrus.WithError(err).Fatal("failed to match inf glob")
	}
	ansAll, err := filepath.Glob(*ansGlob)
	if err != nil {
		logrus.WithError(err).Fatal("failed to match ans glob")
	}
	if len(infAll) != len(ansAll) {
		logrus.Fatal("input and answer files count mismatch")
	}

	program := []string{
		"from hashlib import sha256",
		"from sys import stdin, stdout",
		"from base64 import b64decode",
		"m = sha256()",
		"m.update(stdin.buffer.read())",
		"d = {",
	}
	for i, inf := range infAll {
		ans := ansAll[i]
		logrus.WithFields(logrus.Fields{
			"inf": filepath.Base(inf), "ans": filepath.Base(ans),
		}).Info("calculate hash")

		mPair, err := pepel.GenMagicPair(inf, ans, *zFlag)
		if err != nil {
			logrus.WithError(err).Fatal("failed to generate magic pair")
		}
		program = append(program, fmt.Sprintf("\t%q: %q,", mPair.InfSHA256, mPair.AnsBase64))
	}
	program = append(program,
		"}", "r = b64decode(d[m.hexdigest()])",
	)
	if *zFlag {
		program = append(program,
			"from zlib import decompress", "r = decompress(r)",
		)
	}
	program = append(program, "stdout.buffer.write(r)")

	fmt.Println(strings.Join(program, "\n")) //nolint:forbidigo // Basic functionality.
}
