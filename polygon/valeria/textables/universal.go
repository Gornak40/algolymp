package textables

import (
	"fmt"
	"strings"
)

const (
	UniversalTag = "universal"
)

// Works both in HTML and PDF render.
type UniversalTable struct {
	groups []string
}

var _ Table = &UniversalTable{}

func (t *UniversalTable) addGroupRow(info GroupInfo, comment string) {
	row := fmt.Sprintf("%s & %d & %s & %s \\\\ \\hline",
		info.Group,
		info.Score,
		comment,
		strings.Join(info.Dependencies, ", "),
	)
	t.groups = append(t.groups, row)
}

func (t *UniversalTable) AddGroup0(info GroupInfo) {
	t.addGroupRow(info, "тесты из условия")
}

func (t *UniversalTable) AddGroup(info GroupInfo) {
	t.addGroupRow(info, "")
}

func (t *UniversalTable) AddLastGroup(info GroupInfo) {
	t.addGroupRow(info, "---")
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
