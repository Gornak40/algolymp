package kultq

import (
	"fmt"
	"os"
	"sort"
	"strconv"

	"github.com/abiosoft/ishell/v2"
)

type Config struct {
	DreamPath  string    `json:"dreamPath"`
	StatBounds []float64 `json:"statBounds"`
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

func (e *Engine) comListFunc(c *ishell.Context) {
	for _, name := range e.probNames {
		p := e.probMapa[name]
		c.Println(p.String())
	}
}

func (e *Engine) comRunFunc(c *ishell.Context) {
	choices := c.Checklist(e.probNames, "Which problems you want to check?", nil)
	sprobs := make([]string, 0, len(choices))
	for _, idx := range choices {
		sprobs = append(sprobs, e.probNames[idx])
	}
	e.runProbs(c, sprobs)
}

func (e *Engine) comSliceFunc(c *ishell.Context) {
	choice := c.MultiChoice(e.probNames, "Which problem you want to slice?")
	if choice == -1 {
		return
	}
	stat := stata{}
	if err := stat.read(e.probNames[choice]); err != nil {
		c.Err(err)

		return
	}

	c.ShowPrompt(false)
	defer c.ShowPrompt(true)
	c.Print("Lower bound (%): ")
	lbs := c.ReadLine()
	lb, err := strconv.Atoi(lbs)
	if err != nil {
		c.Err(err)

		return
	}
	lbf := float64(lb) / 100 //nolint:mnd // percents
	c.Print("Upper bound (%): ")
	rbs := c.ReadLine()
	rb, err := strconv.Atoi(rbs)
	if err != nil {
		c.Err(err)

		return
	}
	rbf := float64(rb) / 100 //nolint:mnd // percents

	var slip []statPair
	for _, p := range stat.pairs {
		if lbf <= p.score && p.score <= rbf {
			slip = append(slip, p)
		}
	}

	c.Printf("Slice size: %d pairs\n", len(slip))
}

func (e *Engine) comStatsFunc(c *ishell.Context) {
	choice := c.MultiChoice(e.probNames, "Which csv report you want to display?")
	if choice == -1 {
		return
	}
	stat := stata{}
	if err := stat.read(e.probNames[choice]); err != nil {
		c.Err(err)

		return
	}
	ps := make([]int, len(e.cfg.StatBounds))
	for _, s := range stat.pairs {
		idx := sort.Search(len(e.cfg.StatBounds), func(i int) bool {
			return e.cfg.StatBounds[i] > s.score
		}) - 1
		ps[idx]++
	}
	c.Printf("Stats for problem %s:\n", e.probNames[choice])
	sum := 0
	for i := len(ps) - 1; i >= 0; i-- {
		sum += ps[i]
		c.Printf("≥ %d%%: %d pairs\n", int(e.cfg.StatBounds[i]*100), sum) //nolint:mnd // percents
	}
}

func (e *Engine) runShell() {
	shell := ishell.New()

	shell.AddCmd(&ishell.Cmd{
		Name: "list",
		Help: "list contest problems",
		Func: e.comListFunc,
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "run",
		Help: "run antiplagiarism comparator",
		Func: e.comRunFunc,
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "slice",
		Help: "check random pairs in fixed interval",
		Func: e.comSliceFunc,
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "stats",
		Help: "show csv report percentiles",
		Func: e.comStatsFunc,
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
	bar.Prefix("Progress (do not interrupt): ")
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
	c.Println("Done, csv reports are written in contest directory")
}
