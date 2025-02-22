package gibon

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"slices"

	"github.com/Gornak40/algolymp/polygon"
	"github.com/sirupsen/logrus"
)

const (
	ModeContest  = "contest"
	ModeCommit   = "commit"
	ModeDownload = "download"
	ModeGroups   = "groups"
	ModePackage  = "package"
	ModeUpdate   = "update"
)

const packageMode = 0666

var (
	ErrNoPackage     = errors.New("no suitable package")
	ErrUnknownMethod = errors.New("unknown method")
	ErrBadGroupDesc  = errors.New("bad group description")
)

type Gibon struct {
	client *polygon.Polygon
	pID    int
}

func NewGibon(client *polygon.Polygon, pID int) *Gibon {
	return &Gibon{
		client: client,
		pID:    pID,
	}
}

func (g *Gibon) resolveDownload() error {
	prob, err := g.client.GetProblem(g.pID) // it's for zip naming
	if err != nil {
		logrus.WithError(err).Fatal("failed to get problem")
	}
	logrus.WithFields(logrus.Fields{
		"name": prob.Name, "owner": prob.Owner, "access": prob.AccessType,
		"package": prob.LatestPackage, "revision": prob.Revision,
	}).Info("problem found")

	pkgs, err := g.client.GetPackages(g.pID)
	if err != nil {
		return err
	}
	idx := slices.IndexFunc(pkgs, func(p polygon.PackageAnswer) bool {
		return p.State == "READY" && p.Type == "linux" && p.Revision == prob.Revision
	})
	if idx == -1 {
		return ErrNoPackage
	}
	p := pkgs[idx]
	logrus.WithFields(logrus.Fields{
		"revision": p.Revision, "comment": p.Comment, "type": p.Type,
	}).Info("package found")

	data, err := g.client.DownloadPackage(g.pID, p.ID, p.Type)
	if err != nil {
		return err
	}
	fname := fmt.Sprintf("%s-%d-%s.zip", prob.Name, p.Revision, p.Type)
	logrus.WithField("filename", fname).Info("save package")

	return os.WriteFile(fname, data, packageMode)
}

func (g *Gibon) listProblems() error {
	probs, err := g.client.ContestProblems(g.pID)
	if err != nil {
		return err
	}
	for _, p := range probs {
		fmt.Println(p.ID) //nolint:forbidigo // Basic functionality.
	}

	return nil
}

func (g *Gibon) markGroups() error {
	if err := g.client.EnableGroups(g.pID); err != nil {
		return err
	}
	if err := g.client.EnablePoints(g.pID); err != nil {
		return err
	}
	r := csv.NewReader(os.Stdin)
	logrus.Info("waiting for TODO input...")
	for {
		ln, err := r.Read()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return err
		}
		if len(ln) != 2 { //nolint:mnd // line length
			return fmt.Errorf("%w: invalid line length", ErrBadGroupDesc)
		}
		var l, r int
		if _, err := fmt.Sscanf(ln[1], "%d-%d", &l, &r); err != nil {
			return err
		}
		if r < l {
			return fmt.Errorf("%w: r < l", ErrBadGroupDesc)
		}
		logrus.WithFields(logrus.Fields{"group": ln[0], "l": l, "r": r}).Info("set group")
		ids := make([]int, 0, r-l+1)
		for i := l; i <= r; i++ {
			ids = append(ids, i)
		}
		if err := g.client.SetTestGroup(g.pID, ln[0], ids); err != nil {
			return err
		}
	}

	return nil
}

func (g *Gibon) Resolve(method string) error {
	switch method {
	case ModeContest:
		return g.listProblems()
	case ModeCommit:
		return g.client.Commit(g.pID, true, "")
	case ModeDownload:
		return g.resolveDownload()
	case ModePackage:
		return g.client.BuildPackage(g.pID, true, true)
	case ModeUpdate:
		return g.client.UpdateWorkingCopy(g.pID)
	case ModeGroups:
		return g.markGroups()
	}

	return fmt.Errorf("%w: %s", ErrUnknownMethod, method)
}
