package polygon

import (
	"fmt"
	"net/url"
	"strconv"
)

const (
	defaultTestset = "tests"
)

type TestRequest url.Values

func NewTestRequest(pID int, index int) TestRequest {
	return TestRequest{
		"problemId": []string{strconv.Itoa(pID)},
		"testIndex": []string{strconv.Itoa(index)},
		"testset":   []string{defaultTestset},
	}
}

func (tr TestRequest) TestSet(testset string) TestRequest {
	tr["testset"] = []string{testset}

	return tr
}

func (tr TestRequest) Group(group string) TestRequest {
	tr["testGroup"] = []string{group}

	return tr
}

func (tr TestRequest) Points(points float32) TestRequest {
	tr["testPoints"] = []string{fmt.Sprint(points)}

	return tr
}

func (tr TestRequest) Input(input string) TestRequest {
	tr["testInput"] = []string{input}

	return tr
}

func (tr TestRequest) Description(description string) TestRequest {
	tr["testDescription"] = []string{description}

	return tr
}

func (tr TestRequest) UseInStatements(f bool) TestRequest {
	tr["testUseInStatements"] = []string{strconv.FormatBool(f)}

	return tr
}

type ProblemRequest url.Values

func NewProblemRequest(pID int) ProblemRequest {
	return ProblemRequest{
		"problemId": []string{strconv.Itoa(pID)},
	}
}

func (pr ProblemRequest) InputFile(name string) ProblemRequest {
	pr["inputFile"] = []string{name}

	return pr
}

func (pr ProblemRequest) OutputFile(name string) ProblemRequest {
	pr["outputFile"] = []string{name}

	return pr
}

func (pr ProblemRequest) Interactive(f bool) ProblemRequest {
	pr["interactive"] = []string{strconv.FormatBool(f)}

	return pr
}

func (pr ProblemRequest) TimeLimit(tl int) ProblemRequest {
	pr["timeLimit"] = []string{strconv.Itoa(tl)}

	return pr
}

func (pr ProblemRequest) MemoryLimit(ml int) ProblemRequest {
	pr["memoryLimit"] = []string{strconv.Itoa(ml)}

	return pr
}

type FileRequest url.Values

func NewFileRequest(pID int, typ FileType, name, file string) FileRequest {
	return FileRequest{
		"problemId": []string{strconv.Itoa(pID)},
		"type":      []string{string(typ)},
		"name":      []string{name},
		"file":      []string{file},
	}
}

func (fr FileRequest) CheckExisting(f bool) FileRequest {
	fr["checkExisting"] = []string{strconv.FormatBool(f)}

	return fr
}

func (fr FileRequest) SourceType(typ string) FileRequest {
	fr["sourceType"] = []string{typ}

	return fr
}

// TODO: add other options

type SolutionRequest url.Values

func NewSolutionRequest(pID int, name, file string, tag SolutionTag) SolutionRequest {
	return SolutionRequest{
		"problemId": []string{strconv.Itoa(pID)},
		"name":      []string{name},
		"file":      []string{file},
		"tag":       []string{string(tag)},
	}
}

func (sr SolutionRequest) CheckExisting(f bool) SolutionRequest {
	sr["checkExisting"] = []string{strconv.FormatBool(f)}

	return sr
}

func (sr SolutionRequest) SourceType(typ string) SolutionRequest {
	sr["sourceType"] = []string{typ}

	return sr
}

type StatementRequest url.Values

func NewStatementRequest(pID int, lang string) StatementRequest {
	return StatementRequest{
		"problemId": []string{strconv.Itoa(pID)},
		"lang":      []string{lang},
	}
}

func (sr StatementRequest) Encoding(enc string) StatementRequest {
	sr["encoding"] = []string{enc}

	return sr
}

func (sr StatementRequest) Name(name string) StatementRequest {
	sr["name"] = []string{name}

	return sr
}

func (sr StatementRequest) Legend(legend string) StatementRequest {
	sr["legend"] = []string{legend}

	return sr
}

func (sr StatementRequest) Input(input string) StatementRequest {
	sr["input"] = []string{input}

	return sr
}

func (sr StatementRequest) Output(output string) StatementRequest {
	sr["output"] = []string{output}

	return sr
}

func (sr StatementRequest) Scoring(scoring string) StatementRequest {
	sr["scoring"] = []string{scoring}

	return sr
}

func (sr StatementRequest) Interaction(interaction string) StatementRequest {
	sr["interaction"] = []string{interaction}

	return sr
}

func (sr StatementRequest) Notes(notes string) StatementRequest {
	sr["notes"] = []string{notes}

	return sr
}

func (sr StatementRequest) Tutorial(tutorial string) StatementRequest {
	sr["tutorial"] = []string{tutorial}

	return sr
}

type ValidatorTestRequest url.Values

func NewValidatorTestRequest(pID, index int) ValidatorTestRequest {
	return ValidatorTestRequest{
		"problemId": []string{strconv.Itoa(pID)},
		"testIndex": []string{strconv.Itoa(index)},
	}
}

func (vtr ValidatorTestRequest) Input(input string) ValidatorTestRequest {
	vtr["testInput"] = []string{input}

	return vtr
}

// VALID or INVALID.
func (vtr ValidatorTestRequest) Verdict(verdict string) ValidatorTestRequest {
	vtr["testVerdict"] = []string{verdict}

	return vtr
}

type CheckerTestRequest url.Values

func NewCheckerTestRequest(pID, index int) CheckerTestRequest {
	return CheckerTestRequest{
		"problemId": []string{strconv.Itoa(pID)},
		"testIndex": []string{strconv.Itoa(index)},
	}
}

func (ctr CheckerTestRequest) Input(input string) CheckerTestRequest {
	ctr["testInput"] = []string{input}

	return ctr
}

func (ctr CheckerTestRequest) Answer(answer string) CheckerTestRequest {
	ctr["testAnswer"] = []string{answer}

	return ctr
}

func (ctr CheckerTestRequest) Output(output string) CheckerTestRequest {
	ctr["testOutput"] = []string{output}

	return ctr
}

func (ctr CheckerTestRequest) Verdict(verdict string) CheckerTestRequest {
	ctr["testVerdict"] = []string{verdict}

	return ctr
}
