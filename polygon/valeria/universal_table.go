package valeria

import (
	"fmt"
	"strings"
)

// Works both in HTML and PDF render.
type UniversalTable struct {
	groups []string
}

var _ TexTable = &UniversalTable{}

func (t *UniversalTable) addGroupRow(info groupInfo, comment string) {
	row := fmt.Sprintf("%s & %d & %s & %s \\\\ \\hline",
		info.group,
		info.score,
		comment,
		strings.Join(info.dependencies, ", "),
	)
	t.groups = append(t.groups, row)
}

func (t *UniversalTable) addGroup0(info groupInfo) {
	t.addGroupRow(info, "тесты из условия")
}

func (t *UniversalTable) addGroup(info groupInfo) {
	t.addGroupRow(info, "")
}

func (t *UniversalTable) addLastGroup(info groupInfo) {
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
