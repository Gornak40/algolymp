package miniparse_test

import (
	"strings"
	"testing"

	"github.com/Gornak40/algolymp/pkg/miniparse"
	"github.com/stretchr/testify/require"
)

const coreMini = `[core]
id = avx2024
name = AVX инструкции с нуля 🦍
number_id = 2812
number_id = 1233
`

func TestSimple(t *testing.T) {
	t.Parallel()
	r := strings.NewReader(coreMini)
	var ss struct {
		Core struct {
			ID   string `mini:"id"`
			Name string `mini:"name"`
		} `mini:"core"`
	}

	require.NoError(t, miniparse.Decode(r, &ss))
}
