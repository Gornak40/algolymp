package polygon

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
)

const (
	FullProblemScore = 100
)

var (
	ErrAllTestsAreSamples = fmt.Errorf("all tests are samples, try -s flag")
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
	s := Scoring{
		score:        map[string]int{},
		count:        map[string]int{},
		first:        map[string]int{},
		last:         map[string]int{},
		dependencies: map[string][]string{},
	}
	for _, t := range tests {
		s.score[t.Group] += int(t.Points) // TODO: ensure ejudge doesn't support float points
		s.count[t.Group]++
		if val, ok := s.first[t.Group]; !ok || val > t.Index {
			s.first[t.Group] = t.Index
		}
		if val, ok := s.last[t.Group]; !ok || val < t.Index {
			s.last[t.Group] = t.Index
		}
	}
	for _, g := range groups {
		if g.PointsPolicy != "COMPLETE_GROUP" {
			return nil, errors.New("test_score is not supported yet")
		}
		if s.last[g.Name]-s.first[g.Name]+1 != s.count[g.Name] {
			return nil, errors.New("bad tests order, fix in polygon required")
		}
		s.dependencies[g.Name] = g.Dependencies
		s.groups = append(s.groups, g.Name)
	}
	return &s, nil
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
	for i, g := range s.groups {
		var info string
		if i == 0 {
			info = "тесты из условия"
		} else if i == len(s.groups)-1 {
			info = "---"
		}
		ans = append(ans, fmt.Sprintf("%s & %d & %s & %s \\\\ \\hline",
			g, s.score[g], info, strings.Join(s.dependencies[g], ", ")))
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

	s, err := NewScoring(tests, groups)
	if err != nil {
		return err
	}
	valuer := s.buildValuer()
	if verbose {
		logrus.Info("valuer.cfg\n" + valuer)
	}

	link := p.buildURL("problem.saveFile", url.Values{
		"problemId": []string{strconv.Itoa(pID)},
		"type":      []string{"resource"},
		"name":      []string{"valuer.cfg"},
		"file":      []string{valuer},
	})
	if _, err := p.makeQuery(http.MethodPost, link); err != nil {
		return err
	}

	scoring := s.buildScoring()
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
	tc := 0
	for _, t := range tests {
		if t.UseInStatements && !samples {
			continue
		}
		tc++
	}
	if tc == 0 {
		return ErrAllTestsAreSamples
	}
	small := FullProblemScore / tc
	smallCnt := tc - (FullProblemScore - small*tc)
	logrus.WithFields(logrus.Fields{
		"zeroCount":  len(tests) - tc,
		"smallScore": small,
		"smallCount": smallCnt,
		"bigScore":   small + 1,
		"bigCount":   tc - smallCnt,
	}).Info("points statistics")
	for _, t := range tests {
		var gr string
		var pt int
		if smallCnt == 0 { //nolint:gocritic // It's smart piece of code.
			gr = "2"
			pt = small + 1
		} else if !t.UseInStatements || samples {
			gr = "1"
			pt = small
			smallCnt--
		} else {
			gr = "0"
			pt = 0
		}
		rt := NewTestRequest(pID, t.Index).
			Group(gr).
			Points(float32(pt))
		if err := p.SaveTest(rt); err != nil {
			return err
		}
	}
	return nil
}
