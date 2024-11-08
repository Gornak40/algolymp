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

func TestSimple(t *testing.T) {
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
