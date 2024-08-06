package vydra

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/Gornak40/algolymp/internal/natstream"
	"github.com/Gornak40/algolymp/polygon"
	"github.com/sirupsen/logrus"
)

const (
	megabyte = 1024 * 1024

	defaultTL     = 1000
	defaultML     = 256
	defaultInput  = "input"
	defaultOutput = "output"

	chkTests = "files/tests/checker-tests"
)

var (
	ErrBadSolutionTag = errors.New("bad solution tag")
)

type Vydra struct {
	client    *polygon.Polygon
	pID       int
	prob      ProblemXML
	streamIn  *natstream.NatStream
	streamOut *natstream.NatStream
	streamAns *natstream.NatStream
}

func NewVydra(client *polygon.Polygon, pID int) *Vydra {
	return &Vydra{
		client:    client,
		pID:       pID,
		streamIn:  new(natstream.NatStream),
		streamOut: new(natstream.NatStream),
		streamAns: new(natstream.NatStream),
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

// TODO: add validator tests.
func (v *Vydra) initValidator(val *Validator) error {
	logrus.WithFields(logrus.Fields{
		"path": val.Source.Path, "type": val.Source.Type,
	}).Info("init validator")

	return v.client.SetValidator(v.pID, filepath.Base(val.Source.Path))
}

// TODO: add checker tests.
func (v *Vydra) initChecker(chk *Checker) error {
	path := chk.Name
	if path == "" {
		path = filepath.Base(chk.Source.Path)
	}
	logrus.WithFields(logrus.Fields{
		"path": path, "type": chk.Type,
	}).Info("init checker")

	return v.client.SetChecker(v.pID, path)
}

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

// TODO: add points, groups, etc.
func (v *Vydra) uploadTest(testset string, idx int, test *Test) error {
	logrus.WithFields(logrus.Fields{
		"testset": testset, "idx": idx, "method": test.Method, "sample": test.Sample,
	}).Info("upload test")

	tr := polygon.NewTestRequest(v.pID, idx).
		TestSet(testset).
		Description(test.Description).
		UseInStatements(test.Sample)
	if test.Method == "manual" {
		input, err := v.streamIn.Next()
		if err != nil {
			return err
		}
		tr.Input(input)
	}

	return v.client.SaveTest(tr)
}

func (v *Vydra) uploadScript(testset *TestSet) error {
	logrus.WithField("testset", testset.Name).Info("upload script")
	gens := make([]string, 0, testset.TestCount)
	for idx, test := range testset.Tests.Tests { // build script
		if test.Method == "generated" {
			gens = append(gens, fmt.Sprintf("%s > %d", test.Cmd, idx+1))
		}
	}
	script := strings.Join(gens, "\n")

	return v.client.SaveScript(v.pID, testset.Name, script)
}

func (v *Vydra) uploadValidatorTest(idx int, test *Test) error {
	logrus.WithFields(logrus.Fields{"idx": idx}).Info("upload validator test")
	input, err := v.streamIn.Next()
	if err != nil {
		return err
	}

	vtr := polygon.NewValidatorTestRequest(v.pID, idx).
		Input(input).Verdict(strings.ToUpper(test.Verdict))

	return v.client.SaveValidatorTest(vtr)
}

func (v *Vydra) uploadCheckerTest(idx int, test *Test) error {
	logrus.WithFields(logrus.Fields{"idx": idx}).Info("upload checker test")
	input, err := v.streamIn.Next()
	if err != nil {
		return err
	}
	output, err := v.streamOut.Next()
	if err != nil {
		return err
	}
	answer, err := v.streamAns.Next()
	if err != nil {
		return err
	}

	ctr := polygon.NewCheckerTestRequest(v.pID, idx).
		Input(input).Output(output).Answer(answer).
		Verdict(strings.ReplaceAll(strings.ToUpper(test.Verdict), "-", "_")) // oh my God

	return v.client.SaveCheckerTest(ctr)
}

//nolint:funlen,cyclop // it's good design
func (v *Vydra) Upload(errs chan error) error {
	defer close(errs)
	if err := v.readXML("problem.xml"); err != nil {
		return err
	}
	errs <- v.initProblem(&v.prob.Judging)
	errs <- v.uploadTags(v.prob.Tags.Tags)
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
	if val := v.prob.Assets.Validators.Validator; val != nil {
		errs <- v.initValidator(val)
		if err := v.streamIn.Init("files/tests/validator-tests/*"); err != nil {
			errs <- err

			goto checker
		}
		for idx, test := range val.TestSet.Tests.Tests {
			errs <- v.uploadValidatorTest(idx+1, &test)
		}
	}
checker:
	if chk := v.prob.Assets.Checker; chk != nil {
		errs <- v.initChecker(chk)
		if err := v.streamIn.Init(filepath.Join(chkTests, "*[^.ao]")); err != nil {
			errs <- err

			goto judging
		}
		if err := v.streamOut.Init(filepath.Join(chkTests, "*.o")); err != nil {
			errs <- err

			goto judging
		}
		if err := v.streamAns.Init(filepath.Join(chkTests, "*.a")); err != nil {
			errs <- err

			goto judging
		}
		for idx, test := range chk.TestSet.Tests.Tests {
			errs <- v.uploadCheckerTest(idx+1, &test)
		}
	}
judging:
	for _, testset := range v.prob.Judging.TestSets {
		errs <- v.uploadScript(&testset)
		if err := v.streamIn.Init(path.Join(testset.Name, "*[^.a]")); err != nil {
			errs <- err

			continue
		}
		for idx, test := range testset.Tests.Tests {
			errs <- v.uploadTest(testset.Name, idx+1, &test)
		}
	}

	return nil
}
