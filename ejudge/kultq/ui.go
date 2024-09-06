package kultq

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"

	"github.com/abiosoft/ishell/v2"
)

const idFmt = "%d"

func (e *Engine) addDeathNote(path string) error {
	var id int
	_, _ = fmt.Sscanf(filepath.Base(path), idFmt, &id)
	if _, err := e.deathNotes.Write([]byte(strconv.Itoa(id) + "\n")); err != nil {
		return err
	}

	return nil
}

func (e *Engine) startUI(c *ishell.Context, slip []statPair) error {
	printCode := func(path string) error {
		args := make([]string, len(e.cfg.BatArgs), len(e.cfg.BatArgs)+1)
		copy(args, e.cfg.BatArgs)
		args = append(args, path)
		bat := exec.Command(e.cfg.BatBin, args...) //nolint:gosec // this is the way
		bat.Stdout = os.Stdout

		return bat.Run()
	}

	for _, s := range slip {
		if err := c.ClearScreen(); err != nil {
			return err
		}
		if err := printCode(s.path1); err != nil {
			return err
		}
		if err := printCode(s.path2); err != nil {
			return err
		}
		c.Println("Score (%):", s.score*100) //nolint:mnd //percents
		yn, err := confirm(c, "Kill this guys (Ctrl+D to stop)?")
		if errors.Is(err, io.EOF) {
			return nil
		}
		if err != nil {
			return err
		}
		if yn {
			if err := e.addDeathNote(s.path1); err != nil {
				return err
			}
			if err := e.addDeathNote(s.path2); err != nil {
				return err
			}
		}
	}

	return nil
}
