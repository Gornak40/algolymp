package polygon

import (
	"errors"

	"github.com/sirupsen/logrus"
)

const (
	FullProblemScore = 100
)

var (
	ErrAllTestsAreSamples = errors.New("all tests are samples, try -s flag")
	ErrBadTestsOrder      = errors.New("bad tests order, fix in polygon required")
)

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
		switch {
		case smallCnt == 0:
			group = "2"
			points = small + 1
		case !test.UseInStatements || samples:
			group = "1"
			points = small
			smallCnt--
		default:
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
