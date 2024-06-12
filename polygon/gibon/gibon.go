package gibon

import (
	"errors"
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

var (
	ErrNoPackages = errors.New("no suitable packages")
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

func (g *Gibon) Resolve(method string) error {
	switch method {
	case ModeCommit:
		if err := g.client.Commit(g.pID, true, ""); err != nil {
			return err
		}
	case ModeDownload:
		pkgs, err := g.client.GetPackages(g.pID)
		if err != nil {
			return err
		}
		idx := slices.IndexFunc(pkgs, func(p polygon.PackageAnswer) bool {
			return p.State == "READY" && p.Type == "linux"
		})
		if idx == -1 {
			return ErrNoPackages
		}
		p := pkgs[idx]
		logrus.WithFields(logrus.Fields{
			"revision": p.Revision,
			"comment":  p.Comment,
			"type":     p.Type,
		}).Info("package found")
	case ModePackage:
		if err := g.client.BuildPackage(g.pID, true, true); err != nil {
			return err
		}
	case ModeUpdate:
		if err := g.client.UpdateWorkingCopy(g.pID); err != nil {
			return err
		}
	}

	return nil
}
