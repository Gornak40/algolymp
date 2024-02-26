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
