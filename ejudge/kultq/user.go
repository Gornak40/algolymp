package kultq

import (
	"os"
	"path"
)

type user struct {
	directory string
	pyRuns    []string
	cppRuns   []string
}

func (p *problem) initUser(dir string) (*user, error) {
	usr := user{directory: dir}
	runs, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	for _, run := range runs {
		rname := run.Name()
		ext := path.Ext(rname)
		switch {
		case run.IsDir():
			p.nestedCnt++
		case ext == ".cpp": // C++
			usr.cppRuns = append(usr.cppRuns, path.Join(dir, rname))
		case ext == ".py": // Python
			usr.pyRuns = append(usr.pyRuns, path.Join(dir, rname))
		default:
			p.unknown[ext]++
		}
	}

	return &usr, nil
}
