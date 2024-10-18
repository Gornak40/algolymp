//go:build windows

package printer

import "github.com/sirupsen/logrus"

func PrintFile(name, device string) error {
	logrus.WithFields(logrus.Fields{"file": name, "printer": device}).Info("print file")

	return nil
}
