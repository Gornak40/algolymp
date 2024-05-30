package textables

import (
	"fmt"
	"strings"
)

// If len(vars) != 0 works only in PDF render.
// Otherwise works both in HTML and PDF render.
type MoscowTable struct {
	groups []string
	vars   []string
}

const clineDelta = 2

var _ Table = &MoscowTable{}

func NewMoscowTable(vars []string) *MoscowTable {
	newVars := make([]string, 0, len(vars))
	for _, dv := range vars {
		newVars = append(newVars, "$"+dv+"$")
	}

	return &MoscowTable{
		vars: newVars,
	}
}

func (t *MoscowTable) AddGroup(info GroupInfo) {
	var comment, limits string
	if info.Type == Group0 {
		comment = "Тесты из условия."
	}
	if len(t.vars) > 0 {
		limits = strings.Repeat(" & ", len(t.vars))
	} else {
		limits = " & "
	}
	row := fmt.Sprintf("%s & %d & %s%s & %s \\\\ \\hline",
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
	}
	if len(t.vars) > 0 {
		table = append(table,
			fmt.Sprintf("\\begin{tabular}{|c|c|%sc|c|}", strings.Repeat("c|", len(t.vars))),
			"\\hline",
			fmt.Sprintf("& & \\multicolumn{%d}{c|}{Доп. ограничения} & & \\\\", len(t.vars)),
			fmt.Sprintf("\\cline{3-%d}", len(t.vars)+clineDelta),
			"\\raisebox{2.25ex}[0cm][0cm]{Группа}",
			"& \\raisebox{2.25ex}[0cm][0cm]{Баллы}",
			"& "+strings.Join(t.vars, " & "),
			"& \\raisebox{2.25ex}[0cm][0cm]{Необх. группы}",
			"& \\raisebox{2.25ex}[0cm][0cm]{Комментарий}",
		)
	} else {
		table = append(table,
			"\\begin{tabular}{|c|c|c|c|c|}",
			"\\hline",
			"Группа",
			"& Баллы",
			"& Доп. ограничения",
			"& Необх. группы",
			"& Комментарий",
		)
	}
	table = append(table, "\\\\", "\\hline")
	table = append(table, t.groups...)
	table = append(table, "\\end{tabular}", "\\end{center}")

	return strings.Join(table, "\n")
}
