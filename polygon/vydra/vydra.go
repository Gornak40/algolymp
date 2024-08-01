package vydra

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/Gornak40/algolymp/polygon"
	"github.com/sirupsen/logrus"
)

const (
	statementType = "application/x-tex"
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
	case "time-limit-exceeded-or-accepted":
		tag = polygon.TagTLorOK
	case "presentation-error":
		tag = polygon.TagPresentationError
	case "memory-limit-exceeded":
		tag = polygon.TagMemoryLimit
	default:
		return fmt.Errorf("%w: %s", ErrBadSolutionTag, sol.Tag)
	}

	sr := polygon.NewSolutionRequest(v.pID, filepath.Base(sol.Source.Path), string(data), tag).
		SourceType(sol.Source.Type)

	return v.client.SaveSolution(sr)
}

func (v *Vydra) uploadResource(res *File) error {
	logrus.WithFields(logrus.Fields{
		"path": res.Path, "type": res.Type,
	}).Info("upload resource")
	data, err := os.ReadFile(res.Path)
	if err != nil {
		return err
	}

	fr := polygon.NewFileRequest(v.pID, polygon.TypeResource, filepath.Base(res.Path), string(data))

	return v.client.SaveFile(fr)
}

func (v *Vydra) uploadExecutable(exe *Executable) error {
	logrus.WithFields(logrus.Fields{
		"path": exe.Source.Path, "type": exe.Source.Type,
	}).Info("upload executable")
	data, err := os.ReadFile(exe.Source.Path)
	if err != nil {
		return err
	}

	fr := polygon.NewFileRequest(v.pID, polygon.TypeSource, filepath.Base(exe.Source.Path), string(data)).
		SourceType(exe.Source.Type)

	return v.client.SaveFile(fr)
}

func (v *Vydra) uploadStatement(stat *Statement) error {
	if stat.Type != statementType {
		return nil
	}
	logrus.WithFields(logrus.Fields{
		"language": stat.Language, "type": stat.Type, "charset": stat.Charset,
	}).Info("upload statement")
	dir := "statement-sections/" + stat.Language

	return filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		sr := polygon.NewStatementRequest(v.pID, stat.Language).
			Encoding(stat.Charset)
		switch filepath.Base(path) {
		case "input.tex":
			sr.Input(string(data))
		case "output.tex":
			sr.Output(string(data))
		case "legend.tex":
			sr.Legend(string(data))
		case "name.tex":
			sr.Name(string(data))
		case "notes.tex":
			sr.Notes(string(data))
		case "tutorial.tex":
			sr.Tutorial(string(data))
		case "interaction.tex":
			sr.Interaction(string(data))
		case "scoring.tex":
			sr.Scoring(string(data))
		default:
			return nil
		}
		logrus.WithField("path", path).Info("upload statement section")

		return v.client.SaveStatement(sr)
	})
}

func (v *Vydra) uploadTags(tags []Tag) error {
	stags := make([]string, 0, len(tags))
	for _, t := range tags {
		stags = append(stags, t.Value)
	}
	line := strings.Join(stags, ",")
	logrus.WithField("tags", line).Info("upload tags")

	return v.client.SaveTags(v.pID, line)
}

// TODO: add validator tests.
func (v *Vydra) uploadValidator(val *Validator) error {
	logrus.WithFields(logrus.Fields{
		"path": val.Source.Path, "type": val.Source.Type,
	}).Info("upload validator")

	return v.client.SetValidator(v.pID, filepath.Base(val.Source.Path))
}

// TODO: add checker tests.
func (v *Vydra) uploadChecker(chk *Checker) error {
	path := chk.Name
	if path == "" {
		path = filepath.Base(chk.Source.Path)
	}
	logrus.WithFields(logrus.Fields{
		"path": path, "type": chk.Type,
	}).Info("upload checker")

	return v.client.SetChecker(v.pID, path)
}

func (v *Vydra) Upload(errs chan error) error {
	defer close(errs)
	if err := v.readXML("problem.xml"); err != nil {
		return err
	}
	for _, sol := range v.prob.Assets.Solutions.Solutions {
		errs <- v.uploadSolution(&sol)
	}
	for _, res := range v.prob.Files.Resources.Files {
		errs <- v.uploadResource(&res)
	}
	for _, exe := range v.prob.Files.Executables.Executables {
		errs <- v.uploadExecutable(&exe)
	}
	for _, stat := range v.prob.Statements.Statements {
		errs <- v.uploadStatement(&stat)
	}
	errs <- v.uploadTags(v.prob.Tags.Tags)
	if val := v.prob.Assets.Validators.Validator; val != nil {
		errs <- v.uploadValidator(val)
	}
	if chk := v.prob.Assets.Checker; chk != nil {
		errs <- v.uploadChecker(chk)
	}

	return nil
}
