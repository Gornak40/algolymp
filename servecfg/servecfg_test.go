package servecfg_test

import (
	"strings"
	"testing"

	"github.com/Gornak40/algolymp/servecfg"
	"github.com/stretchr/testify/require"
)

const config1 = ` contest_time   = 0
virtual
# comment
score_system = acm # another comment

  [problem]
time_limit = 5
id = 1
use_stdin
short_name = "A"
use_stdout

[problem ]
short_name = "B"
	id = 2
max_vm_size = 256M
time_limit = 10

[language]
id =  2	
short_name  = "gcc"
long_name = "GCC = the best compiler"

 [ problem] 
use_stdin = 0
 id = 3
time_limit = 3
short_name = "C"
use_stdout = 0
`

func TestString(t *testing.T) {
	t.Parallel()
	cfg := servecfg.New(strings.NewReader(config1))

	require.Equal(t, `# -*- coding: utf-8 -*-

contest_time = 0
virtual
score_system = acm


[language]
id = 2
short_name = "gcc"
long_name = "GCC = the best compiler"


[problem]
time_limit = 5
id = 1
use_stdin
short_name = "A"
use_stdout

[problem]
short_name = "B"
id = 2
max_vm_size = 256M
time_limit = 10

[problem]
use_stdin = 0
id = 3
time_limit = 3
short_name = "C"
use_stdout = 0
`, cfg.String())
}

func TestRoot(t *testing.T) {
	t.Parallel()
	cfg := servecfg.New(strings.NewReader(config1))

	require.Equal(t, []servecfg.Field{
		{Key: "contest_time", Value: "0", SectionIdx: 1},
		{Key: "virtual", SectionIdx: 1},
		{Key: "score_system", Value: "acm", SectionIdx: 1},
	}, cfg.Query("."))

	require.Empty(t, cfg.Query(".compile_dir"))

	require.Equal(t, []servecfg.Field{
		{Key: "contest_time", Value: "0", SectionIdx: 1},
	}, cfg.Query(".contest_time"))

	require.Equal(t, []servecfg.Field{
		{Key: "virtual", SectionIdx: 1},
		{Key: "score_system", Value: "acm", SectionIdx: 1},
	}, cfg.Query(".score_system,virtual,id,short_name"))
}

func TestSection(t *testing.T) {
	t.Parallel()
	cfg := servecfg.New(strings.NewReader(config1))

	require.Equal(t, []servecfg.Field{
		{Key: "id", Value: "2", Section: "language", SectionIdx: 1},
		{Key: "short_name", Value: `"gcc"`, Section: "language", SectionIdx: 1},
		{Key: "long_name", Value: `"GCC = the best compiler"`, Section: "language", SectionIdx: 1},
	}, cfg.Query("@language"))

	require.Empty(t, cfg.Query("@language:2"))

	require.Equal(t, []servecfg.Field{
		{Key: "id", Value: "1", Section: "problem", SectionIdx: 1},
		{Key: "short_name", Value: `"A"`, Section: "problem", SectionIdx: 1},
		{Key: "short_name", Value: `"B"`, Section: "problem", SectionIdx: 2},
		{Key: "id", Value: "2", Section: "problem", SectionIdx: 2},
		{Key: "max_vm_size", Value: "256M", Section: "problem", SectionIdx: 2},
		{Key: "id", Value: "3", Section: "problem", SectionIdx: 3},
		{Key: "short_name", Value: `"C"`, Section: "problem", SectionIdx: 3},
	}, cfg.Query("@problem.id,max_vm_size,short_name"))

	require.Equal(t, []servecfg.Field{
		{Key: "time_limit", Value: "5", Section: "problem", SectionIdx: 1},
		{Key: "use_stdin", Section: "problem", SectionIdx: 1},
		{Key: "use_stdin", Value: "0", Section: "problem", SectionIdx: 3},
		{Key: "time_limit", Value: "3", Section: "problem", SectionIdx: 3},
	}, cfg.Query("@problem:1,3.time_limit,use_stdin,contest_time"))
}

// TODO: add unit tests for Update, Set and Delete.
func TestE2E(t *testing.T) {
	t.Parallel()
	cfg := servecfg.New(strings.NewReader(config1))

	qu := cfg.Query("@problem.use_stdout")
	cfg.Update("1", qu)

	qu = cfg.Query("@language")
	cfg.Update(servecfg.Deleter, qu)

	qu = cfg.Query(".score_system")
	cfg.Update("kirov", qu)

	qu = cfg.Query(".virtual")
	cfg.Update(servecfg.Deleter, qu)

	qu = cfg.Query("@problem:1,2.max_vm_size")
	cfg.Set("max_vm_size", "512M", qu)

	qu = cfg.Query("@problem:2")
	cfg.Set("long_name", `"gorilla-and-horror-tree"`, qu)

	qu = cfg.Query("@problem", "@problem:2")
	cfg.Set("use_stdin", "", qu)

	qu = cfg.Query("@problem:1,2,3,4.time_limit")
	cfg.Update("777", qu)

	qu = cfg.Query("@problem:1")
	cfg.Set("max_vm_size", "64M", qu)

	require.Equal(t, `# -*- coding: utf-8 -*-

contest_time = 0
score_system = kirov


[problem]
time_limit = 777
id = 1
use_stdin
short_name = "A"
use_stdout = 1
max_vm_size = 64M

[problem]
short_name = "B"
id = 2
max_vm_size = 512M
time_limit = 777
long_name = "gorilla-and-horror-tree"
use_stdin

[problem]
use_stdin
id = 3
time_limit = 777
short_name = "C"
use_stdout = 1
`, cfg.String())
}
