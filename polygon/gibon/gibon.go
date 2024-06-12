package gibon

import (
	"errors"
	"fmt"
	"os"
	"slices"

	"github.com/Gornak40/algolymp/polygon"
	"github.com/sirupsen/logrus"
)

const (
	ModeCommit   = "commit"
	ModeDownload = "download"
	ModePackage  = "package"
	ModeUpdate   = "update"
)

const packageMode = 0666

var (
	ErrNoPackage     = errors.New("no suitable package")
	ErrUnknownMethod = errors.New("unknown method")
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

func (g *Gibon) Resolve(method string) error {
	switch method {
	case ModeCommit:
		return g.client.Commit(g.pID, true, "")
	case ModeDownload:
		return g.resolveDownload()
	case ModePackage:
		return g.client.BuildPackage(g.pID, true, true)
	case ModeUpdate:
		return g.client.UpdateWorkingCopy(g.pID)
	}

	return fmt.Errorf("%w: %s", ErrUnknownMethod, method)
}
