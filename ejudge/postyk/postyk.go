package postyk

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"

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
	cfg    *ejudge.Config
	client *http.Client
	target string
}

func NewIndexer(cfg *ejudge.Config) *Indexer {
	return &Indexer{
		cfg:    cfg,
		client: http.DefaultClient,
	}
}

func (i *Indexer) Feed(cID int) error {
	var err error
	i.target, err = url.JoinPath(i.cfg.URL,
		printRoot,
		i.cfg.Secret1,
		fmt.Sprintf("%06d", cID),
		"print")
	if err != nil {
		return fmt.Errorf("failed to build url: %w", err)
	}
	logrus.WithField("url", i.target).Info("init indexer")

	subs, err := i.LoadList()
	if err != nil { // test load
		return err
	}
	logrus.WithField("count", len(subs)).Info("success ping directory")

	return nil
}

func (i *Indexer) GetContent(name string) ([]byte, error) {
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

func (i *Indexer) LoadList() ([]*Submission, error) {
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
