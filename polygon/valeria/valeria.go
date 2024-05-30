package valeria

import (
	"errors"
	"fmt"
	"strings"

	"github.com/Gornak40/algolymp/polygon"
	"github.com/Gornak40/algolymp/polygon/valeria/textables"
	"github.com/sirupsen/logrus"
)

var (
	ErrNoTestScore   = errors.New("test_score is not supported yet")
	ErrBadTestsOrder = errors.New("bad tests order, fix in polygon required")
)

type Valeria struct {
	client *polygon.Polygon
}

func NewValeria(client *polygon.Polygon) *Valeria {
	return &Valeria{
		client: client,
	}
}

type scoring struct {
	score        map[string]int
	count        map[string]int
	first        map[string]int
	last         map[string]int
	dependencies map[string][]string
	groups       []string
}

func newScoring(tests []polygon.TestAnswer, groups []polygon.GroupAnswer) (*scoring, error) {
	scorer := scoring{
		score:        map[string]int{},
		count:        map[string]int{},
		first:        map[string]int{},
		last:         map[string]int{},
		dependencies: map[string][]string{},
	}
	for _, test := range tests {
		scorer.score[test.Group] += int(test.Points)
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

func (s *scoring) buildValuer() string {
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

func (s *scoring) buildScoring(table textables.Table) {
	for index, groupName := range s.groups {
		info := textables.GroupInfo{
			Name:         groupName,
			Score:        s.score[groupName],
			Dependencies: s.dependencies[groupName],
		}
		switch index {
		case 0:
			info.Type = textables.Group0
		case len(s.groups) - 1:
			info.Type = textables.GroupLast
		default:
			info.Type = textables.GroupRegular
		}
		table.AddGroup(info)
	}
}

func (v *Valeria) InformaticsValuer(pID int, table textables.Table, verbose bool) error {
	groups, err := v.client.GetGroups(pID)
	if err != nil {
		return err
	}
	tests, err := v.client.GetTests(pID)
	if err != nil {
		return err
	}

	scorer, err := newScoring(tests, groups)
	if err != nil {
		return err
	}
	valuer := scorer.buildValuer()
	if verbose {
		logrus.Info("valuer.cfg\n" + valuer)
	}
	fr := polygon.NewFileRequest(pID, polygon.TypeResource, "valuer.cfg", valuer)
	if err := v.client.SaveFile(fr); err != nil {
		return err
	}
	scorer.buildScoring(table)

	return nil
}
