//go:build windows

package printer

import (
	"fmt"
	"os/exec"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

func PrintFile(name, device string) error {
	logrus.WithFields(logrus.Fields{"file": filepath.Base(name), "printer": device}).Info("print file")

	command := fmt.Sprintf("Get-Content -Path %q | Out-Printer -Name %q", name, device)
	cmd := exec.Command("powershell", "-Command", command)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%w: %s", err, string(out))
	}
	logrus.Info(string(out))

	return nil
}
