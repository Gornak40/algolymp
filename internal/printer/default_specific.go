//go:build !windows

package printer

import (
	"path/filepath"

	"github.com/sirupsen/logrus"
)

func PrintFile(name, device string) error {
	logrus.WithFields(logrus.Fields{"file": filepath.Base(name), "printer": device}).Info("print file")
	logrus.Warn("not implemented yet (use Windows)")

	return nil
}
