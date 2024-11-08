package miniparse_test

import (
	"strings"
	"testing"

	"github.com/Gornak40/algolymp/pkg/miniparse"
	"github.com/stretchr/testify/require"
)

const coreMini = `[core]
# year is just for fun
id = avx2024
name = AVX инструкции с нуля 🦍

number_id = 2812
number_id = 1233
public = 1
contest_id = 33
magic = -1

flags = false
flags = true
flags = T

contest_id = 12312
contest_id = 9012
`

func TestReflect(t *testing.T) {
	t.Parallel()
	type core struct {
		ID        string `mini:"id"`
		Name      string `mini:"name"`
		NumberID  []int  `mini:"number_id"`
		Public    bool   `mini:"public"`
		Empty     string
		Flags     []bool   `mini:"flags"`
		ContestID []string `mini:"contest_id"`
		Magic     int      `mini:"magic"`
	}
	type config struct {
		Core core `mini:"core"`
	}
	var ss config

	r := strings.NewReader(coreMini)
	require.NoError(t, miniparse.Decode(r, &ss))
	require.Equal(t, config{
		Core: core{
			ID:        "avx2024",
			Name:      "AVX инструкции с нуля 🦍",
			NumberID:  []int{2812, 1233},
			Public:    true,
			Empty:     "",
			Flags:     []bool{false, true, true},
			ContestID: []string{"33", "12312", "9012"},
			Magic:     -1,
		},
	}, ss)
}

func TestCorner(t *testing.T) {
	t.Parallel()
	var ss struct{}
	tt := []string{
		"",
		"[alone]\n",
		"\n",
		"[empty]\n[empty]\n[again]\n",
		"[rune]\nkey = [section]\n[section]\nkey2 = 0\n",
	}

	for _, s := range tt {
		r := strings.NewReader(s)
		require.NoError(t, miniparse.Decode(r, &ss))
	}
}

func TestParseError(t *testing.T) {
	t.Parallel()
	var ss struct{}
	tt := map[string]error{
		"[html page]\n":                            miniparse.ErrInvalidSection,
		"[htmlPage]\n":                             miniparse.ErrInvalidSection,
		"[1html]\n":                                miniparse.ErrInvalidSection,
		"[html]":                                   miniparse.ErrUnexpectedEOF,
		"[html":                                    miniparse.ErrUnexpectedEOF,
		"[html]\n]":                                miniparse.ErrInvalidChar,
		"[html]\nkey\n":                            miniparse.ErrInvalidKey,
		"[html]\nkey=value\n":                      miniparse.ErrInvalidKey,
		"[html]\nKEY = value\n":                    miniparse.ErrInvalidChar,
		"[html]\n1key = value\n":                   miniparse.ErrInvalidChar,
		"[html]\nkey= value\n":                     miniparse.ErrInvalidKey,
		"[html]\nkey =value\n":                     miniparse.ErrExpectedSpace,
		"[html]\nkey = value":                      miniparse.ErrUnexpectedEOF,
		"[html]\nkey  = value\n":                   miniparse.ErrExpectedEqual,
		"[html]\nkey = value[html]\n[html]":        miniparse.ErrUnexpectedEOF,
		"[html] # section\n":                       miniparse.ErrExpectedNewLine,
		"# section\n[html]\nkey = 1\nkey2\n":       miniparse.ErrInvalidKey,
		"[html]\nkey = \"bababababayka\nbrrry\"\n": miniparse.ErrInvalidKey,
		" [html]\n":                                miniparse.ErrLeadingSpace,
		"[html]\n key = value\n":                   miniparse.ErrLeadingSpace,
		"[html]\n = value\n":                       miniparse.ErrLeadingSpace,
		"key = value\n":                            miniparse.ErrRootSection,
	}

	for s, terr := range tt {
		r := strings.NewReader(s)
		require.ErrorIs(t, miniparse.Decode(r, &ss), terr)
	}
}

func TestReflectArgError(t *testing.T) {
	t.Parallel()

	m := make(map[string]any)
	require.ErrorIs(t, miniparse.Decode(strings.NewReader(""), m), miniparse.ErrExpectedPointer)
	require.ErrorIs(t, miniparse.Decode(strings.NewReader(""), &m), miniparse.ErrExpectedStruct)
	require.ErrorIs(t, miniparse.Decode(strings.NewReader(""), nil), miniparse.ErrExpectedPointer)

	x := 4
	require.ErrorIs(t, miniparse.Decode(strings.NewReader(""), &x), miniparse.ErrExpectedStruct)
	require.ErrorIs(t, miniparse.Decode(strings.NewReader(""), "gorill"), miniparse.ErrExpectedPointer)

	var ss struct{}
	require.ErrorIs(t, miniparse.Decode(strings.NewReader(""), ss), miniparse.ErrExpectedPointer)
}

const htmlMini = `[html]
page_id = 14141
preload = true
# wow! a comment
page_id = 99119
name = Bob
`

func TestReflectError(t *testing.T) {
	t.Parallel()
	r := strings.NewReader(htmlMini)

	var ss2 struct {
		HTML struct {
			PageID []struct {
				Key   string `mini:"key"`
				Value string `mini:"value"`
			} `mini:"page_id"`
		} `mini:"html"`
	}
	require.ErrorIs(t, miniparse.Decode(r, &ss2), miniparse.ErrBadRecordType)

	var ss3 struct {
		HTML struct {
			Name   string `mini:"name"`
			PageID int    `mini:"page_id"`
		} `mini:"html"`
	}
	r.Reset(htmlMini)
	require.ErrorIs(t, miniparse.Decode(r, &ss3), miniparse.ErrExpectedArray)
}
