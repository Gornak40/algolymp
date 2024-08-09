package vydra

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/Gornak40/algolymp/polygon"
	"github.com/sirupsen/logrus"
)

func (v *Vydra) initProblem(judge *Judging) error {
	input := defaultInput
	if judge.InputFile != "" {
		input = judge.InputFile
	}
	output := defaultOutput
	if judge.OutputFile != "" {
		output = judge.OutputFile
	}
	tl := defaultTL
	ml := defaultML
	if len(judge.TestSets) != 0 {
		tl = judge.TestSets[0].TimeLimit
		ml = judge.TestSets[0].MemoryLimit / megabyte
	}
	logrus.WithFields(logrus.Fields{
		"input": input, "output": output,
		"tl": tl, "ml": ml,
	}).Info("init problem")

	pr := polygon.NewProblemRequest(v.pID).
		InputFile(input).OutputFile(output).
		TimeLimit(tl).MemoryLimit(ml)

	return v.client.UpdateInfo(pr)
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

func (v *Vydra) uploadStatement(stat *Statement) error {
	if stat.Type != "application/x-tex" {
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

func (v *Vydra) batchInitial(errs chan error) {
	errs <- v.initProblem(&v.prob.Judging)
	if tags := v.prob.Tags.Tags; len(tags) != 0 {
		errs <- v.uploadTags(tags)
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
}
