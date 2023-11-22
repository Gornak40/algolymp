package polygon

import (
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
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
	cfg *Config
}

func NewPolygon(cfg *Config) *Polygon {
	return &Polygon{
		cfg: cfg,
	}
}

func (p *Polygon) makeQuery(method string, params url.Values) ([]byte, error) {
	url, _ := url.JoinPath(p.cfg.URL, "api", method)
	logrus.WithFields(logrus.Fields{"params": params.Encode()}).Info(method)

	params["apiKey"] = []string{p.cfg.APIKey}
	params["time"] = []string{fmt.Sprint(time.Now().Unix())}
	sig := fmt.Sprintf("%s/%s?%s#%s", sixSecretSymbols, method, params.Encode(), p.cfg.APISecret)

	b := sha512.Sum512([]byte(sig))
	hsh := hex.EncodeToString(b[:])
	params["apiSig"] = []string{sixSecretSymbols + hsh}

	url = fmt.Sprintf("%s?%s", url, params.Encode())
	resp, err := http.Get(url) //nolint:gosec,noctx // it's just get query, relax
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
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
	Tests          int      `json:"tests"`
}

type Answer struct {
	Status  string          `json:"status"`
	Comment string          `json:"comment"`
	Result  json.RawMessage `json:"result"`
}

func (p *Polygon) getGroups(pID int) ([]GroupAnswer, error) {
	data, err := p.makeQuery("problem.viewTestGroup", url.Values{
		"problemId": []string{fmt.Sprint(pID)},
		"testset":   []string{"tests"},
	})
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
	data, err := p.makeQuery("problem.tests", url.Values{
		"problemId": []string{fmt.Sprint(pID)},
		"testset":   []string{"tests"},
	})
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

func (p *Polygon) GetValuer(pID int) error {
	groups, err := p.getGroups(pID)
	if err != nil {
		return err
	}

	tests, err := p.getTests(pID)
	if err != nil {
		return err
	}

	score := map[string]int{}
	count := map[string]int{}
	first := map[string]int{}
	last := map[string]int{}

	for _, t := range tests {
		score[t.Group] += int(t.Points) // TODO: ensure ejudge doesn't support float points
		count[t.Group]++
		if val, ok := first[t.Group]; !ok || val > t.Index {
			first[t.Group] = t.Index
		}
		if val, ok := last[t.Group]; !ok || val < t.Index {
			last[t.Group] = t.Index
		}
	}

	res := []string{}
	for _, g := range groups {
		if g.PointsPolicy != "COMPLETE_GROUP" {
			return errors.New("test_score not supported yet")
		}
		if last[g.Name]-first[g.Name]+1 != count[g.Name] {
			return errors.New("bad tests order, fix in polygon required")
		}
		cur := fmt.Sprintf("group %s {\n\ttests %d-%d;\n\tscore %d;\n",
			g.Name, first[g.Name], last[g.Name], score[g.Name])
		if len(g.Dependencies) != 0 {
			cur += fmt.Sprintf("\trequires %s;\n", strings.Join(g.Dependencies, ","))
		}
		cur += "}\n"
		res = append(res, cur)
	}

	valuer := strings.Join(res, "\n")
	logrus.Info(valuer)

	return nil
}
