package textables

import "github.com/sirupsen/logrus"

type GroupInfo struct {
	Group        string
	Score        int
	Dependencies []string
}

type Table interface {
	AddGroup0(info GroupInfo)
	AddGroup(info GroupInfo)
	AddLastGroup(info GroupInfo)

	String() string
}

func GetTexTable(tableTyp string, cntVars int) Table { //nolint:ireturn
	logrus.WithFields(logrus.Fields{"type": tableTyp, "vars": cntVars}).Info("select textable")

	switch tableTyp {
	case UniversalTag:
		return &UniversalTable{}
	default:
		return nil
	}
}
