package kultq

import (
	"errors"
	"io"
	"os"
	"os/exec"

	"github.com/abiosoft/ishell/v2"
)

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
			c.Println("KILLED") // TODO: implement
		}
	}

	return nil
}
