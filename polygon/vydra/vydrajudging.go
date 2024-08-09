package vydra

import (
	"fmt"
	"path"
	"strings"

	"github.com/Gornak40/algolymp/polygon"
	"github.com/sirupsen/logrus"
)

func (v *Vydra) uploadScript(testset *TestSet) error {
	logrus.WithField("testset", testset.Name).Info("upload script")
	gens := make([]string, 0, testset.TestCount)
	for idx, test := range testset.Tests.Tests { // build script
		if test.Method == "generated" {
			line := fmt.Sprintf("%s > %d", test.Cmd, idx+1)
			if test.FromFile != "" { // TODO: find better solution for `gen > {3-100}`
				line = fmt.Sprintf("%s > {%d}", test.Cmd, idx+1)
			}
			gens = append(gens, line)
		}
	}
	script := strings.Join(gens, "\n")
	if script == "" {
		return nil
	}

	return v.client.SaveScript(v.pID, testset.Name, script)
}

func (v *Vydra) uploadTest(testset string, idx int, test *Test) error {
	// It's kind of experimental solution.
	if (*test == Test{Cmd: test.Cmd, FromFile: test.FromFile, Method: "generated"}) {
		return nil
	}
	logrus.WithFields(logrus.Fields{
		"testset": testset, "idx": idx,
		"method": test.Method, "sample": test.Sample,
	}).Info("upload test")

	tr := polygon.NewTestRequest(v.pID, idx).
		TestSet(testset).
		Description(test.Description).
		UseInStatements(test.Sample)
	if test.Group != "" {
		tr.Group(test.Group)
	}
	if test.Points != 0 {
		tr.Points(test.Points)
	}
	if test.Method == "manual" {
		input, err := v.streamIn.Next()
		if err != nil {
			return err
		}
		tr.Input(input)
	}

	return v.client.SaveTest(tr)
}

func (v *Vydra) initGroups() error {
	logrus.Info("init test groups")

	return v.client.EnableGroups(v.pID)
}

func (v *Vydra) initPoints() error {
	logrus.Info("init test points")

	return v.client.EnablePoints(v.pID)
}

func (v *Vydra) uploadGroup(testset string, group *Group) error {
	deps := make([]string, 0, len(group.Dependencies.Dependencies))
	for _, d := range group.Dependencies.Dependencies {
		deps = append(deps, d.Group)
	}
	logrus.WithFields(logrus.Fields{
		"feedback":     group.FeedbackPolicy,
		"points":       group.PointsPolicy,
		"dependencies": deps,
	}).Info("upload group")

	tgr := polygon.NewTestGroupRequest(v.pID, testset, group.Name).
		FeedbackPolicy(convertString(group.FeedbackPolicy)).
		PointsPolicy(convertString(group.PointsPolicy)).
		Dependencies(deps)

	return v.client.SaveTestGroup(tgr)
}

type testsMetaInfo struct {
	enablePoints bool
	enableGroups bool
}

func getTestsMeta(tests []Test) testsMetaInfo {
	var ans testsMetaInfo
	for _, t := range tests {
		if t.Group != "" {
			ans.enableGroups = true
		}
		if t.Points != 0 {
			ans.enablePoints = true
		}
	}

	return ans
}

func (v *Vydra) batchJudging(errs chan error) {
	for _, testset := range v.prob.Judging.TestSets {
		errs <- v.uploadScript(&testset)
		if err := v.streamIn.Init(path.Join(testset.Name, "*[^.a]")); err != nil {
			errs <- err

			continue
		}
		meta := getTestsMeta(testset.Tests.Tests)
		if meta.enableGroups {
			errs <- v.initGroups()
		}
		if meta.enablePoints {
			errs <- v.initPoints()
		}
		for idx, test := range testset.Tests.Tests {
			errs <- v.uploadTest(testset.Name, idx+1, &test)
		}
		if grp := testset.Groups.Groups; len(grp) != 0 {
			for _, g := range grp {
				errs <- v.uploadGroup(testset.Name, &g)
			}
		}
	}
}
