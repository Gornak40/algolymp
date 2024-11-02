package textables

import (
	"fmt"
	"strings"
)

// Works both in HTML and PDF render.
type KalugaTable struct {
	groups []string
}

var _ Table = &KalugaTable{}

func (t *KalugaTable) AddGroup(info GroupInfo) {
	var limits string
	switch info.Type {
	case Group0:
		return
	case GroupLast:
		limits = "Ограничения из условия"
	case GroupRegular:
	}
	row := fmt.Sprintf("%s & %d & %s \\\\ \\hline",
		info.Name,
		info.Score,
		limits,
	)
	t.groups = append(t.groups, row)
}

func (t *KalugaTable) String() string {
	table := []string{
		"\\begin{center}",
		"\\begin{tabular}{|c|c|c|}",
		"\\hline",
		"\\bf{Подзадача} &",
		"\\bf{Баллы} &",
		"\\bf{Ограничения}",
		"\\\\ \\hline",
	}
	table = append(table, t.groups...)
	table = append(table, "\\end{tabular}", "\\end{center}")

	return strings.Join(table, "\n")
}
