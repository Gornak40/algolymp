package polygon

import (
	"fmt"
	"net/url"
	"strconv"
)

type TestRequest url.Values

func NewTestRequest(pID int, index int) TestRequest {
	return TestRequest{
		"problemId": []string{strconv.Itoa(pID)},
		"testIndex": []string{strconv.Itoa(index)},
		"testset":   []string{defaultTestset},
	}
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

type FileRequest url.Values

func NewFileRequest(pID int, typ FileType, name, file string) FileRequest {
	return FileRequest{
		"problemId": []string{strconv.Itoa(pID)},
		"type":      []string{string(typ)},
		"name":      []string{name},
		"file":      []string{file},
	}
}

// TODO: fix it.
func (fr FileRequest) CheckExisting(f bool) FileRequest {
	fr["checkExisting"] = []string{strconv.FormatBool(f)}

	return fr
}

func (fr FileRequest) SourceType(typ string) FileRequest {
	fr["sourceType"] = []string{typ}

	return fr
}

// TODO: add other options
