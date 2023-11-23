package polygon

import (
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/sirupsen/logrus"
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
			return nil, errors.New("test_score not supported yet")
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

func (p *Polygon) InformaticsValuer(pID int) error {
	groups, err := p.getGroups(pID)
	if err != nil {
		return err
	}
	tests, err := p.getTests(pID)
	if err != nil {
		return err
	}

	s, err := NewScoring(tests, groups)
	if err != nil {
		return err
	}
	valuer := s.buildValuer()
	logrus.Info(valuer)

	link := p.buildURL("problem.saveFile", url.Values{
		"problemId": []string{fmt.Sprint(pID)},
		"type":      []string{"resource"},
		"name":      []string{"valuer.cfg"},
		"file":      []string{valuer},
	})
	if err := p.postQuery(link); err != nil {
		return err
	}

	scoring := s.buildScoring()
	logrus.Info(scoring)

	link = p.buildURL("problem.saveStatement", url.Values{
		"problemId": []string{fmt.Sprint(pID)},
		"lang":      []string{"russian"},
		"scoring":   []string{scoring},
	})
	return p.postQuery(link)
}
