package vydra

import (
	"encoding/xml"
	"errors"
	"fmt"
	"os"

	"github.com/Gornak40/algolymp/polygon"
	"github.com/sirupsen/logrus"
)

var (
	ErrBadSolutionTag = errors.New("bad solution tag")
)

type Vydra struct {
	client *polygon.Polygon
	pID    int
	prob   ProblemXML
}

func NewVydra(client *polygon.Polygon, pID int) *Vydra {
	return &Vydra{
		client: client,
		pID:    pID,
	}
}

func (v *Vydra) readXML(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	if err := xml.Unmarshal(data, &v.prob); err != nil {
		return err
	}
	logrus.WithFields(logrus.Fields{
		"revision": v.prob.Revision, "short-name": v.prob.ShortName,
	}).Info("load problem.xml")

	return nil
}

func (v *Vydra) uploadSolution(sol *Solution) error {
	logrus.WithFields(logrus.Fields{
		"path": sol.Source.Path, "type": sol.Source.Type, "tag": sol.Tag,
	}).Info("upload solution")
	data, err := os.ReadFile(sol.Source.Path)
	if err != nil {
		return err
	}

	var tag polygon.SolutionTag
	switch sol.Tag { // TODO: add other tags
	case "main":
		tag = polygon.TagMain
	case "accepted":
		tag = polygon.TagCorrect
	case "rejected":
		tag = polygon.TagIncorrect
	case "time-limit-exceeded":
		tag = polygon.TagTimeLimit
	case "wrong-answer":
		tag = polygon.TagWrongAnswer
	default:
		return fmt.Errorf("%w: %s", ErrBadSolutionTag, sol.Tag)
	}

	sr := polygon.NewSolutionRequest(v.pID, sol.Source.Path, string(data), tag).
		SourceType(sol.Source.Type)

	return v.client.SaveSolution(sr)
}

func (v *Vydra) Upload() error {
	if err := v.readXML("problem.xml"); err != nil {
		return err
	}
	for _, sol := range v.prob.Assets.Solutions.Solutions {
		if err := v.uploadSolution(&sol); err != nil {
			return err
		}
	}

	return nil
}
