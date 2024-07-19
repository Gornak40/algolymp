package kultq

import (
	"encoding/csv"
	"errors"
	"io"
	"os"
	"sort"
	"strconv"
)

var (
	ErrBadStata = errors.New("bad csv stata file")
)

type statPair struct {
	score float64
	path1 string
	path2 string
}

type stata struct {
	pairs []statPair
}

func (s *stata) read(name string) error {
	file, err := os.Open(name + ".csv")
	if err != nil {
		return err
	}
	defer file.Close()

	r := csv.NewReader(file)
	for {
		line, err := r.Read()
		if errors.Is(err, io.EOF) {
			break
		}
		if len(line) != 3 { //nolint:mnd // csv format
			return ErrBadStata
		}
		score, err := strconv.ParseFloat(line[0], 64)
		if err != nil {
			return err
		}
		if _, err := os.Stat(line[1]); err != nil {
			return err
		}
		if _, err := os.Stat(line[2]); err != nil {
			return err
		}
		p := statPair{score: score, path1: line[1], path2: line[2]}
		s.pairs = append(s.pairs, p)
	}
	sort.SliceStable(s.pairs, func(i, j int) bool {
		return s.pairs[i].score > s.pairs[j].score
	})

	return nil
}
