package polygon

import (
	"errors"
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
)

const (
	FullProblemScore = 100
)

var (
	ErrAllTestsAreSamples = errors.New("all tests are samples, try -s flag")
	ErrNoTestScore        = errors.New("test_score is not supported yet")
	ErrBadTestsOrder      = errors.New("bad tests order, fix in polygon required")
)

type Scoring struct {
	score        map[string]int
	count        map[string]int
	first        map[string]int
	last         map[string]int
	dependencies map[string][]string
	groups       []string
}

func NewScoring(tests []TestAnswer, groups []GroupAnswer) (*Scoring, error) {
	scorer := Scoring{
		score:        map[string]int{},
		count:        map[string]int{},
		first:        map[string]int{},
		last:         map[string]int{},
		dependencies: map[string][]string{},
	}
	for _, test := range tests {
		scorer.score[test.Group] += int(test.Points) // TODO: ensure ejudge doesn't support float points
		scorer.count[test.Group]++
		if val, ok := scorer.first[test.Group]; !ok || val > test.Index {
			scorer.first[test.Group] = test.Index
		}
		if val, ok := scorer.last[test.Group]; !ok || val < test.Index {
			scorer.last[test.Group] = test.Index
		}
	}
	for _, group := range groups {
		if group.PointsPolicy != "COMPLETE_GROUP" {
			return nil, ErrNoTestScore
		}
		if scorer.last[group.Name]-scorer.first[group.Name]+1 != scorer.count[group.Name] {
			return nil, ErrBadTestsOrder
		}
		scorer.dependencies[group.Name] = group.Dependencies
		scorer.groups = append(scorer.groups, group.Name)
	}

	return &scorer, nil
}

func (s *Scoring) buildValuer() string {
	res := []string{}
	for _, g := range s.groups {
		cur := fmt.Sprintf("group %s {\n\ttests %d-%d;\n\tscore %d;\n",
			g, s.first[g], s.last[g], s.score[g])
		if len(s.dependencies[g]) != 0 {
			cur += fmt.Sprintf("\trequires %s;\n", strings.Join(s.dependencies[g], ","))
		}
		cur += "}\n"
		res = append(res, cur)
	}

	return strings.Join(res, "\n")
}

func (s *Scoring) buildScoring() string {
	ans := []string{
		"\\begin{center}",
		"\\begin{tabular}{|c|c|c|c|}",
		"\\hline",
		"\\textbf{Подзадача} &",
		"\\textbf{Баллы} &",
		"\\textbf{Дополнительные ограничения} &",
		"\\textbf{Необходимые подзадачи}",
		"\\\\ \\hline",
	}
	for i, group := range s.groups {
		var info string
		if i == 0 {
			info = "тесты из условия"
		} else if i == len(s.groups)-1 {
			info = "---"
		}
		ans = append(ans, fmt.Sprintf("%s & %d & %s & %s \\\\ \\hline",
			group, s.score[group], info, strings.Join(s.dependencies[group], ", ")))
	}
	ans = append(ans, "\\end{tabular}", "\\end{center}")

	return strings.Join(ans, "\n")
}

func (p *Polygon) InformaticsValuer(pID int, verbose bool) error {
	groups, err := p.GetGroups(pID)
	if err != nil {
		return err
	}
	tests, err := p.GetTests(pID)
	if err != nil {
		return err
	}

	scorer, err := NewScoring(tests, groups)
	if err != nil {
		return err
	}
	valuer := scorer.buildValuer()
	if verbose {
		logrus.Info("valuer.cfg\n" + valuer)
	}
	fr := NewFileRequest(pID, TypeResource, "valuer.cfg", valuer)
	if err := p.SaveFile(fr); err != nil {
		return err
	}

	scoring := scorer.buildScoring()
	fmt.Println(scoring) //nolint:forbidigo // Basic functionality.

	return nil
}

func (p *Polygon) IncrementalScoring(pID int, samples bool) error {
	if err := p.EnablePoints(pID); err != nil {
		return err
	}
	if err := p.EnableGroups(pID); err != nil {
		return err
	}
	tests, err := p.GetTests(pID)
	if err != nil {
		return err
	}
	testsCount := 0
	for _, t := range tests {
		if t.UseInStatements && !samples {
			continue
		}
		testsCount++
	}
	if testsCount == 0 {
		return ErrAllTestsAreSamples
	}
	small := FullProblemScore / testsCount
	smallCnt := testsCount - (FullProblemScore - small*testsCount)
	logrus.WithFields(logrus.Fields{
		"zeroCount":  len(tests) - testsCount,
		"smallScore": small,
		"smallCount": smallCnt,
		"bigScore":   small + 1,
		"bigCount":   testsCount - smallCnt,
	}).Info("points statistics")
	for _, test := range tests {
		var group string
		var points int
		if smallCnt == 0 { //nolint:gocritic // It's smart piece of code.
			group = "2"
			points = small + 1
		} else if !test.UseInStatements || samples {
			group = "1"
			points = small
			smallCnt--
		} else {
			group = "0"
			points = 0
		}
		rt := NewTestRequest(pID, test.Index).
			Group(group).
			Points(float32(points))
		if err := p.SaveTest(rt); err != nil {
			return err
		}
	}

	return nil
}
