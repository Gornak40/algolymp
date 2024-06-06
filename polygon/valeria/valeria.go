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
	ErrUnknownPointsPolicy = errors.New("unknown points policy")
	ErrMissedGroup         = errors.New("missed group, no tests")
	ErrBadTestScore        = errors.New("bad test score, scores different in one group")
)

type Valeria struct {
	client *polygon.Polygon
}

func NewValeria(client *polygon.Polygon) *Valeria {
	return &Valeria{
		client: client,
	}
}

type pointsPolicy int

const (
	policyCompleteGroup pointsPolicy = iota
	policyEachTest
)

type group struct {
	name         string
	score        int
	count        int
	policy       pointsPolicy
	firstIdx     int
	lastIdx      int
	minScore     int
	maxScore     int
	dependencies []string
}

type scoring struct {
	groups []group
}

func newScoring(tests []polygon.TestAnswer, groups []polygon.GroupAnswer) (*scoring, error) {
	grMapa := make(map[string]group)
	for _, test := range tests {
		score := int(test.Points)
		if g, ok := grMapa[test.Group]; !ok {
			grMapa[test.Group] = group{
				name:     test.Group,
				score:    score,
				count:    1,
				firstIdx: test.Index,
				lastIdx:  test.Index,
				minScore: score,
				maxScore: score,
			}
		} else {
			g.score += score
			g.count++
			g.firstIdx = min(g.firstIdx, test.Index)
			g.lastIdx = max(g.lastIdx, test.Index)
			g.minScore = min(g.minScore, score)
			g.maxScore = max(g.maxScore, score)
			grMapa[test.Group] = g
		}
	}
	scorer := scoring{
		groups: make([]group, 0, len(groups)),
	}
	for _, group := range groups {
		g, ok := grMapa[group.Name]
		if !ok {
			return nil, fmt.Errorf("%w: group %s", ErrMissedGroup, group.Name)
		}
		if g.lastIdx-g.firstIdx+1 != g.count {
			return nil, fmt.Errorf("%w: group %s", polygon.ErrBadTestsOrder, group.Name)
		}
		switch group.PointsPolicy {
		case "COMPLETE_GROUP":
			g.policy = policyCompleteGroup
		case "EACH_TEST":
			g.policy = policyEachTest
			if g.minScore != g.maxScore {
				return nil, fmt.Errorf("%w: group %s", ErrBadTestScore, group.Name)
			}
		default:
			return nil, ErrUnknownPointsPolicy
		}
		g.dependencies = group.Dependencies
		scorer.groups = append(scorer.groups, g)
	}

	return &scorer, nil
}

func (s *scoring) buildValuer() string {
	res := []string{}
	for _, g := range s.groups {
		cur := fmt.Sprintf("group %s {\n\ttests %d-%d;\n\t",
			g.name, g.firstIdx, g.lastIdx)
		switch g.policy {
		case policyCompleteGroup:
			cur += fmt.Sprintf("score %d;\n", g.score)
		case policyEachTest:
			cur += fmt.Sprintf("test_score %d;\n", g.minScore)
		}
		if len(g.dependencies) != 0 {
			cur += fmt.Sprintf("\trequires %s;\n", strings.Join(g.dependencies, ","))
		}
		cur += "}\n"
		res = append(res, cur)
	}

	return strings.Join(res, "\n")
}

func (s *scoring) buildScoring(table textables.Table) {
	for index, g := range s.groups {
		info := textables.GroupInfo{
			Name:         g.name,
			Score:        g.score,
			Dependencies: g.dependencies,
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
