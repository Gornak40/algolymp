package postyk

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

const (
	subArgCount = 5
)

var (
	ErrInvalidSubmission = errors.New("invalid submission string")
)

type Submission struct {
	raw      string
	Time     time.Time
	Printer  string
	Location string
	RunID    int
	Name     string
}

func parseSubmission(sub string) (*Submission, error) {
	args := strings.SplitN(sub, "_", subArgCount)
	if len(args) != subArgCount {
		return nil, ErrInvalidSubmission
	}
	utm, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return nil, err
	}
	id, err := strconv.Atoi(args[3])
	if err != nil {
		return nil, err
	}

	return &Submission{
		raw:      sub,
		Time:     time.Unix(utm, 0),
		Printer:  args[1],
		Location: args[2],
		RunID:    id,
		Name:     args[4],
	}, nil
}
