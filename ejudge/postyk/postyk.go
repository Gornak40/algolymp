package postyk

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"

	"github.com/Gornak40/algolymp/ejudge"
	"github.com/PuerkitoBio/goquery"
	"github.com/sirupsen/logrus"
)

const (
	printRoot = "print"
)

var (
	ErrBadStatusCode = errors.New("bad status code")
)

type Indexer struct {
	cfg       *ejudge.Config
	client    *http.Client
	target    string
	cachePath string
}

func NewIndexer(cfg *ejudge.Config) *Indexer {
	return &Indexer{
		cfg:    cfg,
		client: http.DefaultClient,
	}
}

func (i *Indexer) Feed(cID int) error {
	scID := fmt.Sprintf("%06d", cID)

	var err error
	i.cachePath, err = os.UserCacheDir()
	if err != nil {
		return err
	}
	i.cachePath = path.Join(i.cachePath, "algolymp", "postyk", scID)

	i.target, err = url.JoinPath(i.cfg.URL,
		printRoot,
		i.cfg.Secret1,
		scID,
		"print")
	if err != nil {
		return err
	}
	logrus.WithField("url", i.target).Info("init indexer")

	subs, err := i.GetList()
	if err != nil { // initial healthcheck
		return err
	}
	logrus.WithField("count", len(subs)).Info("success ping shared directory")

	return nil
}

func (i *Indexer) GetFile(name string) ([]byte, error) {
	logrus.WithField("name", name).Info("downloading submission")
	link, err := url.JoinPath(i.target, name)
	if err != nil {
		return nil, err
	}
	resp, err := i.client.Get(link) //nolint:noctx // don't need context here
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w (%s)", ErrBadStatusCode, resp.Status)
	}

	return io.ReadAll(resp.Body)
}

func (i *Indexer) GetList() ([]*Submission, error) {
	logrus.Info("refreshing index")
	resp, err := i.client.Get(i.target) //nolint:noctx // don't need context here
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w (%s)", ErrBadStatusCode, resp.Status)
	}
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}
	els := doc.Find("a").Map(func(_ int, s *goquery.Selection) string {
		return s.AttrOr("href", "")
	})

	subs := make([]*Submission, 0, len(els))
	for _, el := range els {
		if el == "../" {
			continue
		}
		sub, err := parseSubmission(el)
		if err != nil { // TODO: maybe just log?
			return nil, err
		}
		subs = append(subs, sub)
	}

	return subs, nil
}