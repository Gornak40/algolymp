package polygon

import (
	"context"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	sixSecretSymbols = "gorill"
)

type Config struct {
	URL       string `json:"url"`
	APIKey    string `json:"apiKey"`
	APISecret string `json:"apiSecret"`
}

type Polygon struct {
	cfg    *Config
	client *http.Client
}

func NewPolygon(cfg *Config) *Polygon {
	return &Polygon{
		cfg:    cfg,
		client: http.DefaultClient,
	}
}

func (p *Polygon) getQuery(url string) ([]byte, error) {
	req, _ := http.NewRequestWithContext(context.TODO(), http.MethodGet, url, nil)
	resp, err := p.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

func (p *Polygon) postQuery(url string) error {
	req, _ := http.NewRequestWithContext(context.TODO(), http.MethodPost, url, nil)
	resp, err := p.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	data, _ := io.ReadAll(resp.Body)
	logrus.Info(string(data))
	return nil
}

func (p *Polygon) skipEscape(params url.Values) string {
	pairs := []string{}
	for k, vals := range params {
		for _, v := range vals {
			pairs = append(pairs, fmt.Sprintf("%s=%s", k, v))
		}
	}
	sort.Strings(pairs)
	return strings.Join(pairs, "&")
}

func (p *Polygon) buildURL(method string, params url.Values) string {
	url, _ := url.JoinPath(p.cfg.URL, "api", method)
	logrus.Info(method)

	params["apiKey"] = []string{p.cfg.APIKey}
	params["time"] = []string{fmt.Sprint(time.Now().Unix())}
	sig := fmt.Sprintf("%s/%s?%s#%s", sixSecretSymbols, method, p.skipEscape(params), p.cfg.APISecret)

	b := sha512.Sum512([]byte(sig))
	hsh := hex.EncodeToString(b[:])
	params["apiSig"] = []string{sixSecretSymbols + hsh}

	return fmt.Sprintf("%s?%s", url, params.Encode())
}

type TestAnswer struct {
	Index  int     `json:"index"`
	Group  string  `json:"group"`
	Points float32 `json:"points"`
}

type GroupAnswer struct {
	Name           string   `json:"name"`
	PointsPolicy   string   `json:"pointsPolicy"`
	FeedbackPolicy string   `json:"feedbackPolicy"`
	Dependencies   []string `json:"dependencies"`
}

type Answer struct {
	Status  string          `json:"status"`
	Comment string          `json:"comment"`
	Result  json.RawMessage `json:"result"`
}

func (p *Polygon) getGroups(pID int) ([]GroupAnswer, error) {
	data, err := p.getQuery(p.buildURL("problem.viewTestGroup", url.Values{
		"problemId": []string{fmt.Sprint(pID)},
		"testset":   []string{"tests"},
	}))
	if err != nil {
		return nil, err
	}
	var ansG Answer
	if err := json.Unmarshal(data, &ansG); err != nil {
		return nil, err
	}
	if ansG.Status != "OK" {
		return nil, errors.New(ansG.Comment)
	}
	var groups []GroupAnswer
	_ = json.Unmarshal(ansG.Result, &groups)
	return groups, nil
}

func (p *Polygon) getTests(pID int) ([]TestAnswer, error) {
	data, err := p.getQuery(p.buildURL("problem.tests", url.Values{
		"problemId": []string{fmt.Sprint(pID)},
		"testset":   []string{"tests"},
	}))
	if err != nil {
		return nil, err
	}
	var ansT Answer
	if err := json.Unmarshal(data, &ansT); err != nil {
		return nil, err
	}
	if ansT.Status != "OK" {
		return nil, errors.New(ansT.Comment)
	}
	var tests []TestAnswer
	_ = json.Unmarshal(ansT.Result, &tests)
	return tests, nil
}
