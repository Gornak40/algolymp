package textables

import (
	"fmt"
	"strings"
)

// Works both in HTML and PDF render.
type UniversalTable struct {
	groups []string
}

var _ Table = &UniversalTable{}

func (t *UniversalTable) AddGroup(info GroupInfo) {
	var limits string
	switch info.Type {
	case Group0:
		limits = "тесты из условия"
	case GroupLast:
		limits = "---"
	case GroupRegular:
	}
	row := fmt.Sprintf("%s & %d & %s & %s \\\\ \\hline",
		info.Name,
		info.Score,
		limits,
		strings.Join(info.Dependencies, ", "),
	)
	t.groups = append(t.groups, row)
}

func (t *UniversalTable) String() string {
	table := []string{
		"\\begin{center}",
		"\\begin{tabular}{|c|c|c|c|}",
		"\\hline",
		"\\textbf{Подзадача} &",
		"\\textbf{Баллы} &",
		"\\textbf{Дополнительные ограничения} &",
		"\\textbf{Необходимые подзадачи}",
		"\\\\ \\hline",
	}
	table = append(table, t.groups...)
	table = append(table, "\\end{tabular}", "\\end{center}")

	return strings.Join(table, "\n")
}
