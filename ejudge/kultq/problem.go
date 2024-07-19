package kultq

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
)

type problem struct {
	name      string
	cfg       *Config // it's for dream path
	users     []user
	writer    *csv.Writer
	nestedCnt int
	unknown   map[string]int
}

func (p *problem) String() string {
	return fmt.Sprintf("Problem %s: %d users; unknown langs: %s",
		p.name, len(p.users), fmt.Sprint(p.unknown))
}

func (p *problem) init() error {
	p.unknown = make(map[string]int)
	dusr, err := os.ReadDir(p.name)
	if err != nil {
		return err
	}
	p.users = make([]user, 0, len(dusr))
	for _, entry := range dusr {
		if !entry.IsDir() {
			continue
		}
		usr, err := p.initUser(path.Join(p.name, entry.Name()))
		if err != nil {
			return err
		}
		p.users = append(p.users, *usr)
	}

	return nil
}

func (p *problem) runDream(lang string, runs1, runs2 []string) error {
	for i := range runs1 {
		for j := i + 1; j < len(runs2); j++ {
			var out bytes.Buffer
			cmd := exec.Command(p.cfg.DreamPath, lang, runs1[i], runs2[j]) //nolint:gosec
			cmd.Stdout = &out
			if err := cmd.Run(); err != nil {
				return err
			}
			output := strings.TrimSpace(out.String())
			_, err := strconv.ParseFloat(output, 64)
			if err != nil {
				return err
			}
			record := []string{output, runs1[i], runs2[j]}
			if err := p.writer.Write(record); err != nil {
				return err
			}
		}
		p.writer.Flush()
	}

	return nil
}

func (p *problem) compareUsr(a, b user) error {
	if err := p.runDream("cpp", a.cppRuns, b.cppRuns); err != nil {
		return err
	}
	if err := p.runDream("python", a.pyRuns, b.pyRuns); err != nil {
		return err
	}

	return nil
}

func (p *problem) checkProblem(prog chan struct{}, errs chan error) {
	fcsv, err := os.Create(p.name + ".csv")
	if err != nil {
		errs <- err

		return
	}
	defer fcsv.Close()
	p.writer = csv.NewWriter(fcsv)

	for i := range p.users {
		for j := i + 1; j < len(p.users); j++ {
			if err := p.compareUsr(p.users[i], p.users[j]); err != nil {
				errs <- err

				return
			}
		}
		prog <- struct{}{}
	}

	errs <- nil
}
