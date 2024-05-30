package textables

import (
	"fmt"
	"strings"
)

// If cntVars != 0 works only in PDF render.
// Otherwise works both in HTML and PDF render.
type MoscowTable struct {
	groups  []string
	cntVars int
}

var _ Table = &MoscowTable{}

func NewMoscowTable(cntVars int) *MoscowTable {
	return &MoscowTable{
		cntVars: cntVars,
	}
}

func (t *MoscowTable) AddGroup(info GroupInfo) {
	var comment, limits string
	if info.Type == Group0 {
		comment = "Тесты из условия."
	}
	row := fmt.Sprintf("%s & %d & %s & %s & %s \\\\ \\hline",
		info.Name,
		info.Score,
		limits,
		strings.Join(info.Dependencies, ", "),
		comment,
	)
	t.groups = append(t.groups, row)
}

func (t *MoscowTable) String() string {
	table := []string{
		"\\begin{center}",
		"\\renewcommand{\\arraystretch}{1.5}",
		"\\begin{tabular}{|c|c|c|c|c|}",
		"\\hline",
		"Группа",
		"& Баллы",
		"& Доп. ограничения",
		"& Необх. группы",
		"& Комментарий",
		"\\\\",
		"\\hline",
	}
	table = append(table, t.groups...)
	table = append(table, "\\end{tabular}", "\\end{center}")

	return strings.Join(table, "\n")
}
