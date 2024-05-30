package servecfg

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"slices"
	"sort"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
)

const Deleter = "-"

type Config struct {
	Fields []Field
}

type Field struct {
	Key        string
	Value      string
	Section    string
	SectionIdx int
}

func (f *Field) String() string {
	var sPath string
	if f.Section != "" {
		sPath = fmt.Sprintf("@%s:%d", f.Section, f.SectionIdx)
	}
	kv := f.Key
	if f.Value != "" {
		kv = fmt.Sprintf("%s = %s", kv, f.Value)
	}

	return fmt.Sprintf("%s.%s", sPath, kv)
}

// This function does not validate serve.cfg.
func New(reader io.Reader) *Config {
	var section string
	counter := map[string]int{
		"": 1,
	}
	var fields []Field

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()
		if idx := strings.Index(line, "#"); idx != -1 {
			line = line[:idx]
		}
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if line[0] == '[' && line[len(line)-1] == ']' { // section declaration
			section = line[1 : len(line)-1]
			counter[section]++

			continue
		}
		var key, value string
		if kv := strings.SplitN(line, "=", 2); len(kv) == 1 { //nolint:mnd // 2 is two
			key = kv[0]
		} else {
			key, value = strings.TrimSpace(kv[0]), strings.TrimSpace(kv[1])
		}
		fields = append(fields, Field{
			Key:        key,
			Value:      value,
			Section:    section,
			SectionIdx: counter[section],
		})
	}
	if err := scanner.Err(); err != nil {
		logrus.WithError(err).Error("bad serve.cfg reader")
	}

	return &Config{Fields: fields}
}

func (c *Config) String() string {
	mapa := make(map[string][]Field)
	for _, f := range c.Fields {
		mapa[f.Section] = append(mapa[f.Section], f)
	}
	sections := make([]string, 0, len(mapa))
	for s := range mapa {
		sections = append(sections, s)
	}
	sort.Strings(sections)

	result := []string{"# -*- coding: utf-8 -*-\n"}
	for _, s := range sections { // mapa[s] is sorted by sectionIdx
		prevIdx := -1
		for _, field := range mapa[s] {
			if field.SectionIdx != prevIdx {
				if field.Section != "" {
					result = append(result, fmt.Sprintf("\n[%s]", field.Section))
				}
				prevIdx = field.SectionIdx
			}
			if field.Value == "" {
				result = append(result, field.Key)
			} else {
				result = append(result, fmt.Sprintf("%s = %s", field.Key, field.Value))
			}
		}
		result = append(result, "")
	}

	return strings.Join(result, "\n")
}

/*
Example queries:

	.
	.score_system
	.virtual
	@language
	@problem:3
	@problem.ignore_prev_ac
	@problem:1.time_limit
	@problem.use_ac_not_ok,ignore_prev_ac
	@problem:3;4;5.time_limit,id
*/
func (c *Config) Query(queries ...string) []Field {
	var result []Field

	for _, query := range queries {
		var pattern struct {
			keys       []string
			section    string
			sectionIDs []int
		}
		if s := regexp.MustCompile(`\.[\w,]+$`).FindString(query); s != "" {
			pattern.keys = strings.Split(s[1:], ",")
		}
		if s := regexp.MustCompile(`^@\w+`).FindString(query); s != "" {
			pattern.section = s[1:]
		}
		if s := regexp.MustCompile(`:[\d,]+`).FindString(query); s != "" {
			for _, sn := range strings.Split(s[1:], ",") {
				if idx, err := strconv.Atoi(sn); err == nil {
					pattern.sectionIDs = append(pattern.sectionIDs, idx)
				}
			}
		}

		logrus.WithFields(logrus.Fields{
			"fields":     pattern.keys,
			"section":    pattern.section,
			"sectionIds": pattern.sectionIDs,
		}).Infof("parse query \"%s\"", query)
		for _, field := range c.Fields {
			if len(pattern.keys) != 0 && !slices.Contains(pattern.keys, field.Key) {
				continue
			}
			if pattern.section != field.Section {
				continue
			}
			if len(pattern.sectionIDs) != 0 && !slices.Contains(pattern.sectionIDs, field.SectionIdx) {
				continue
			}
			result = append(result, field)
		}
	}

	return result
}

func getMatchedMapa(matched []Field) map[Field]struct{} {
	mapa := make(map[Field]struct{})
	for _, f := range matched {
		mapa[f] = struct{}{}
	}

	return mapa
}

type sectionMatch struct {
	name string
	idx  int
}

func (c *Config) Set(key, value string, matched []Field) *Config {
	mapa := getMatchedMapa(matched)
	secMapa := make(map[sectionMatch]struct{})
	for f := range mapa {
		secMapa[sectionMatch{f.Section, f.SectionIdx}] = struct{}{}
	}

	nwf := make([]Field, 0, len(c.Fields))
	for i, f := range c.Fields {
		sec := sectionMatch{f.Section, f.SectionIdx}
		if _, ok := secMapa[sec]; !ok {
			nwf = append(nwf, f)

			continue
		}
		if f.Key == key {
			delete(secMapa, sec)
			f.Value = value
			nwf = append(nwf, f)

			continue
		}
		nwf = append(nwf, f)
		if i+1 == len(c.Fields) || c.Fields[i+1].Section != f.Section || c.Fields[i+1].SectionIdx != f.SectionIdx {
			nwf = append(nwf, Field{
				Key:        key,
				Value:      value,
				Section:    f.Section,
				SectionIdx: f.SectionIdx,
			})
		}
	}

	c.Fields = nwf

	return c
}

func (c *Config) Update(value string, matched []Field) *Config {
	mapa := getMatchedMapa(matched)
	nwf := make([]Field, 0, len(c.Fields))
	for _, field := range c.Fields {
		if _, ok := mapa[field]; !ok {
			goto writeField
		}
		if value == Deleter {
			continue
		}
		field.Value = value
	writeField:
		nwf = append(nwf, field)
	}

	c.Fields = nwf

	return c
}
