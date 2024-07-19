package kultq

import (
	"fmt"
	"os"

	"github.com/abiosoft/ishell/v2"
)

type Config struct {
	OneCheckPairs int     `json:"oneCheckPairs"`
	DreamPath     string  `json:"dreamPath"`
	LeftScore     float64 `json:"lowerScore"`
}

type Engine struct {
	cfg       *Config
	probNames []string
	probMapa  map[string]problem
}

func NewEngine(cfg *Config) *Engine {
	return &Engine{
		cfg:      cfg,
		probMapa: make(map[string]problem),
	}
}

func (e *Engine) Run() error {
	fi, err := os.ReadDir(".")
	if err != nil {
		return err
	}
	for _, f := range fi {
		if !f.IsDir() {
			continue
		}
		e.probNames = append(e.probNames, f.Name())
		p := problem{name: f.Name(), cfg: e.cfg}
		if p.init() != nil {
			return err
		}
		e.probMapa[f.Name()] = p
	}
	e.runShell()

	return nil
}

func (e *Engine) runShell() {
	shell := ishell.New()

	shell.AddCmd(&ishell.Cmd{
		Name: "list",
		Help: "list contest problems",
		Func: func(c *ishell.Context) {
			for _, name := range e.probNames {
				p := e.probMapa[name]
				c.Println(p.String())
			}
		},
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "run",
		Help: "run antiplagiarism comparator",
		Func: func(c *ishell.Context) {
			choices := c.Checklist(e.probNames, "Which problems you want to check?", nil)
			sprobs := make([]string, 0, len(choices))
			for _, idx := range choices {
				sprobs = append(sprobs, e.probNames[idx])
			}
			e.runProbs(c, sprobs)
		},
	})

	shell.Run()
}

func (e *Engine) runProbs(c *ishell.Context, probs []string) {
	totUsr := 0
	prog := make(chan struct{})
	errs := make(chan error)
	for _, name := range probs {
		p := e.probMapa[name]
		c.Println(p.String())
		totUsr += len(p.users)
		go p.checkProblem(prog, errs)
	}

	progUsr := 0
	bar := c.ProgressBar()
	bar.Prefix("progress (do not interrupt): ")
	bar.Start()
	for waitCnt := len(probs); waitCnt > 0; {
		select {
		case <-prog:
			progUsr++
			perc := progUsr * 100 / totUsr //nolint:mnd // percents
			bar.Suffix(fmt.Sprint(" ", perc, "%"))
			bar.Progress(perc)
		case err := <-errs:
			if err != nil {
				c.Err(err)
			}
			waitCnt--
		}
	}
	bar.Stop()
	c.Println("done, csv reports are written in contest directory")
}
