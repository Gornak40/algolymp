package kultq

import (
	"fmt"
	"os"
	"strings"

	"github.com/abiosoft/ishell/v2"
	"github.com/fatih/color"
)

type Config struct {
	OneCheckPairs int     `json:"oneCheckPairs"`
	DreamPath     string  `json:"dreamPath"`
	LeftScore     float64 `json:"lowerScore"`
}

type Engine struct {
	cfg      *Config
	problems []problem
}

func NewEngine(cfg *Config) *Engine {
	return &Engine{
		cfg: cfg,
	}
}

func (e *Engine) GetProblems() []string {
	probs := make([]string, 0, len(e.problems))
	for _, p := range e.problems {
		probs = append(probs, p.name)
	}

	return probs
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
		p := problem{name: f.Name(), cfg: e.cfg}
		if p.init() != nil {
			return err
		}
		e.problems = append(e.problems, p)
	}
	e.runShell()

	return nil
}

func (e *Engine) runShell() {
	shell := ishell.New()
	probs := e.GetProblems()

	shell.AddCmd(&ishell.Cmd{
		Name: "list",
		Help: "list contest problems",
		Func: func(c *ishell.Context) {
			list := strings.Join(probs, " ")
			c.Println(color.YellowString(list))
		},
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "run",
		Help: "run antiplagiarism comparator",
		Func: func(c *ishell.Context) {
			choices := c.Checklist(probs, "Which problems you want to check?", nil)
			sprobs := make([]string, 0, len(choices))
			for _, idx := range choices {
				sprobs = append(sprobs, probs[idx])
			}
			e.runProbs(c, sprobs)
		},
	})

	shell.Run()
}

func (e *Engine) runProbs(c *ishell.Context, probs []string) {
	pMapa := make(map[string]struct{})
	for _, name := range probs {
		pMapa[name] = struct{}{}
	}
	totUsr := 0
	prog := make(chan struct{})
	errs := make(chan error)
	for _, p := range e.problems {
		if _, ok := pMapa[p.name]; ok {
			c.Printf("problem %s: %d users, %d nested directories, %v unknown languages\n",
				p.name, len(p.users), p.nestedCnt, p.unknown)
			totUsr += len(p.users)
			go p.checkProblem(prog, errs)
		}
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
