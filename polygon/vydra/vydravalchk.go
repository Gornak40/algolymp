package vydra

import (
	"path/filepath"

	"github.com/Gornak40/algolymp/polygon"
	"github.com/sirupsen/logrus"
)

func (v *Vydra) initValidator(val *Validator) error {
	logrus.WithFields(logrus.Fields{
		"path": val.Source.Path, "type": val.Source.Type,
	}).Info("init validator")

	return v.client.SetValidator(v.pID, filepath.Base(val.Source.Path))
}

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

func (v *Vydra) uploadValidatorTest(idx int, test *Test) error {
	logrus.WithFields(logrus.Fields{"idx": idx}).Info("upload validator test")
	input, err := v.streamIn.Next()
	if err != nil {
		return err
	}

	vtr := polygon.NewValidatorTestRequest(v.pID, idx).
		Input(input).Verdict(convertString(test.Verdict))

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
		Verdict(convertString(test.Verdict))

	return v.client.SaveCheckerTest(ctr)
}

func (v *Vydra) batchValChk(errs chan error) {
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

			return
		}
		if err := v.streamOut.Init(filepath.Join(chkTests, "*.o")); err != nil {
			errs <- err

			return
		}
		if err := v.streamAns.Init(filepath.Join(chkTests, "*.a")); err != nil {
			errs <- err

			return
		}
		for idx, test := range chk.TestSet.Tests.Tests {
			errs <- v.uploadCheckerTest(idx+1, &test)
		}
	}
}
