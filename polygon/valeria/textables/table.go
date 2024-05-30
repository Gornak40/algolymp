package textables

type GroupType int

const (
	Group0 GroupType = iota
	GroupRegular
	GroupLast
)

type GroupInfo struct {
	Name         string
	Score        int
	Dependencies []string
	Type         GroupType
}

type Table interface {
	AddGroup(info GroupInfo)
	String() string
}
