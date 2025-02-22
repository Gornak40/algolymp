package polygon

import (
	"bytes"
	"context"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	sixSecretSymbols = "gorill"
)

type SolutionTag string

const (
	TagMain              SolutionTag = "MA"
	TagCorrect           SolutionTag = "OK"
	TagIncorrect         SolutionTag = "RJ"
	TagTimeLimit         SolutionTag = "TL"
	TagTLorOK            SolutionTag = "TO"
	TagWrongAnswer       SolutionTag = "WA"
	TagPresentationError SolutionTag = "PE"
	TagMemoryLimit       SolutionTag = "ML"
	TagRuntimeError      SolutionTag = "RE"
)

type FileType string

const (
	TypeSource   FileType = "source"
	TypeResource FileType = "resource"
	TypeAUX      FileType = "aux"
)

var (
	ErrBadPolygonStatus = errors.New("bad polygon status")
	ErrInvalidMethod    = errors.New("invalid method")
	ErrProblemNotFound  = errors.New("problem not found")
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
	logrus.WithField("url", cfg.URL).Info("init polygon engine")

	return &Polygon{
		cfg:    cfg,
		client: http.DefaultClient,
	}
}

func buildRequest(method, link string, params url.Values) (*http.Request, error) {
	logrus.WithFields(logrus.Fields{
		"method": method,
		"url":    link,
	}).Info("build request")

	switch method {
	case http.MethodGet:
		link = fmt.Sprintf("%s?%s", link, params.Encode())

		return http.NewRequestWithContext(context.TODO(), method, link, nil)
	case http.MethodPost:
		buf := &bytes.Buffer{}
		writer := multipart.NewWriter(buf)
		for k, vals := range params {
			for _, v := range vals {
				wr, err := writer.CreateFormFile(k, k)
				if err != nil {
					return nil, err
				}
				if _, err := wr.Write([]byte(v)); err != nil {
					return nil, err
				}
			}
		}
		if err := writer.Close(); err != nil {
			return nil, err
		}
		req, err := http.NewRequestWithContext(context.TODO(), method, link, buf)
		if err != nil {
			return nil, err
		}
		req.Header.Set("Content-Type", writer.FormDataContentType())

		return req, nil
	default:
		return nil, ErrInvalidMethod
	}
}

func (p *Polygon) makeQuery(method, link string, params url.Values) (*Answer, error) {
	req, err := buildRequest(method, link, params)
	if err != nil {
		return nil, err
	}
	resp, err := p.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var ans Answer
	if err := json.Unmarshal(data, &ans); err != nil {
		return nil, err
	}
	if ans.Status != "OK" {
		return nil, fmt.Errorf("%w: %s", ErrBadPolygonStatus, ans.Comment)
	}

	return &ans, nil
}

func (p *Polygon) skipEscape(params url.Values) string {
	type pair struct {
		key   string
		value string
	}

	var pairs []pair
	for k, vals := range params {
		for _, v := range vals {
			pairs = append(pairs, pair{key: k, value: v})
		}
	}
	sort.Slice(pairs, func(i, j int) bool {
		if pairs[i].key != pairs[j].key {
			return pairs[i].key < pairs[j].key
		}

		return pairs[i].value < pairs[j].value
	})

	pairs2 := make([]string, 0, len(pairs))
	for _, p := range pairs {
		pairs2 = append(pairs2, fmt.Sprintf("%s=%s", p.key, p.value))
	}

	return strings.Join(pairs2, "&")
}

func (p *Polygon) buildURL(method string, params url.Values) (string, url.Values) {
	url, _ := url.JoinPath(p.cfg.URL, "api", method)

	params.Set("apiKey", p.cfg.APIKey)
	params.Set("time", strconv.FormatInt(time.Now().Unix(), 10))
	sig := fmt.Sprintf("%s/%s?%s#%s", sixSecretSymbols, method, p.skipEscape(params), p.cfg.APISecret)

	b := sha512.Sum512([]byte(sig))
	hsh := hex.EncodeToString(b[:])
	params.Set("apiSig", sixSecretSymbols+hsh)

	return url, params
}

func (p *Polygon) BuildPackage(pID int, full, verify bool) error {
	link, params := p.buildURL("problem.buildPackage", url.Values{
		"problemId": {strconv.Itoa(pID)},
		"full":      {strconv.FormatBool(full)},
		"verify":    {strconv.FormatBool(verify)},
	})
	_, err := p.makeQuery(http.MethodPost, link, params)

	return err
}

// Problem Idx (A, B, C) -> Problem.
func (p *Polygon) ContestProblems(pID int) (map[string]ProblemAnswer, error) {
	link, params := p.buildURL("contest.problems", url.Values{
		"contestId": {strconv.Itoa(pID)},
	})
	ansC, err := p.makeQuery(http.MethodGet, link, params)
	if err != nil {
		return nil, err
	}
	var problems map[string]ProblemAnswer
	if err := json.Unmarshal(ansC.Result, &problems); err != nil {
		return nil, err
	}

	return problems, nil
}

func (p *Polygon) Commit(pID int, minor bool, message string) error {
	link, params := p.buildURL("problem.commitChanges", url.Values{
		"problemId":    {strconv.Itoa(pID)},
		"minorChanges": {strconv.FormatBool(minor)},
		"message":      {message},
	})
	_, err := p.makeQuery(http.MethodPost, link, params)

	return err
}

func (p *Polygon) UpdateWorkingCopy(pid int) error {
	link, params := p.buildURL("problem.updateWorkingCopy", url.Values{
		"problemId": {strconv.Itoa(pid)},
	})
	_, err := p.makeQuery(http.MethodPost, link, params)

	return err
}

func (p *Polygon) GetPackages(pID int) ([]PackageAnswer, error) {
	link, params := p.buildURL("problem.packages", url.Values{
		"problemId": {strconv.Itoa(pID)},
	})
	ansP, err := p.makeQuery(http.MethodGet, link, params)
	if err != nil {
		return nil, err
	}
	var packages []PackageAnswer
	if err := json.Unmarshal(ansP.Result, &packages); err != nil {
		return nil, err
	}

	return packages, nil
}

func (p *Polygon) GetGroups(pID int) ([]GroupAnswer, error) {
	link, params := p.buildURL("problem.viewTestGroup", url.Values{
		"problemId": {strconv.Itoa(pID)},
		"testset":   {defaultTestset},
	})
	ansG, err := p.makeQuery(http.MethodGet, link, params)
	if err != nil {
		return nil, err
	}
	var groups []GroupAnswer
	if err := json.Unmarshal(ansG.Result, &groups); err != nil {
		return nil, err
	}

	return groups, nil
}

func (p *Polygon) GetProblem(pID int) (*ProblemAnswer, error) {
	link, params := p.buildURL("problems.list", url.Values{
		"id": {strconv.Itoa(pID)},
	})
	ansP, err := p.makeQuery(http.MethodGet, link, params)
	if err != nil {
		return nil, err
	}
	var problems []ProblemAnswer
	if err := json.Unmarshal(ansP.Result, &problems); err != nil {
		return nil, err
	}
	if len(problems) == 0 {
		return nil, ErrProblemNotFound
	}

	return &problems[0], nil
}

func (p *Polygon) GetTests(pID int) ([]TestAnswer, error) {
	link, params := p.buildURL("problem.tests", url.Values{
		"problemId": {strconv.Itoa(pID)},
		"testset":   {defaultTestset},
		"noInputs":  {"true"},
	})
	ansT, err := p.makeQuery(http.MethodGet, link, params)
	if err != nil {
		return nil, err
	}
	var tests []TestAnswer
	if err := json.Unmarshal(ansT.Result, &tests); err != nil {
		return nil, err
	}

	return tests, nil
}

func (p *Polygon) DownloadPackage(pID, packID int, packType string) ([]byte, error) {
	link, params := p.buildURL("problem.package", url.Values{
		"problemId": {strconv.Itoa(pID)},
		"packageId": {strconv.Itoa(packID)},
		"type":      {packType},
	})
	req, err := buildRequest(http.MethodPost, link, params)
	if err != nil {
		return nil, err
	}
	resp, err := p.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

func (p *Polygon) EnableGroups(pID int) error {
	link, params := p.buildURL("problem.enableGroups", url.Values{
		"problemId": {strconv.Itoa(pID)},
		"testset":   {defaultTestset},
		"enable":    {"true"},
	})
	_, err := p.makeQuery(http.MethodPost, link, params)

	return err
}

func (p *Polygon) EnablePoints(pID int) error {
	link, params := p.buildURL("problem.enablePoints", url.Values{
		"problemId": {strconv.Itoa(pID)},
		"enable":    {"true"},
	})
	_, err := p.makeQuery(http.MethodPost, link, params)

	return err
}

func (p *Polygon) SetTestGroup(pID int, group string, tests []int) error {
	st := make([]string, 0, len(tests))
	for _, v := range tests {
		st = append(st, strconv.Itoa(v))
	}
	link, params := p.buildURL("problem.setTestGroup", url.Values{
		"problemId":   {strconv.Itoa(pID)},
		"testset":     {defaultTestset},
		"testGroup":   {group},
		"testIndices": {strings.Join(st, ",")},
	})
	_, err := p.makeQuery(http.MethodPost, link, params)

	return err
}

func (p *Polygon) SaveFile(fReq FileRequest) error {
	link, params := p.buildURL("problem.saveFile", url.Values(fReq))
	_, err := p.makeQuery(http.MethodPost, link, params)

	return err
}

func (p *Polygon) SaveTest(tReq TestRequest) error {
	link, params := p.buildURL("problem.saveTest", url.Values(tReq))
	_, err := p.makeQuery(http.MethodPost, link, params)

	return err
}

func (p *Polygon) SaveTags(pID int, tags string) error {
	link, params := p.buildURL("problem.saveTags", url.Values{
		"problemId": {strconv.Itoa(pID)},
		"tags":      {tags},
	})
	_, err := p.makeQuery(http.MethodPost, link, params)

	return err
}

func (p *Polygon) SetValidator(pID int, validator string) error {
	link, params := p.buildURL("problem.setValidator", url.Values{
		"problemId": {strconv.Itoa(pID)},
		"validator": {validator},
	})
	_, err := p.makeQuery(http.MethodPost, link, params)

	return err
}

func (p *Polygon) SetChecker(pID int, checker string) error {
	link, params := p.buildURL("problem.setChecker", url.Values{
		"problemId": {strconv.Itoa(pID)},
		"checker":   {checker},
	})
	_, err := p.makeQuery(http.MethodPost, link, params)

	return err
}

func (p *Polygon) UpdateInfo(pr ProblemRequest) error {
	link, params := p.buildURL("problem.updateInfo", url.Values(pr))
	_, err := p.makeQuery(http.MethodPost, link, params)

	return err
}

func (p *Polygon) SetInteractor(pID int, interactor string) error {
	link, params := p.buildURL("problem.setInteractor", url.Values{
		"problemId":  {strconv.Itoa(pID)},
		"interactor": {interactor},
	})
	_, err := p.makeQuery(http.MethodPost, link, params)

	return err
}

func (p *Polygon) SaveScript(pID int, testset, source string) error {
	link, params := p.buildURL("problem.saveScript", url.Values{
		"problemId": {strconv.Itoa(pID)},
		"testset":   {testset},
		"source":    {source},
	})
	_, err := p.makeQuery(http.MethodPost, link, params)

	return err
}

func (p *Polygon) SaveSolution(sr SolutionRequest) error {
	link, params := p.buildURL("problem.saveSolution", url.Values(sr))
	_, err := p.makeQuery(http.MethodPost, link, params)

	return err
}

func (p *Polygon) SaveStatement(sr StatementRequest) error {
	link, params := p.buildURL("problem.saveStatement", url.Values(sr))
	_, err := p.makeQuery(http.MethodPost, link, params)

	return err
}

func (p *Polygon) SaveValidatorTest(vtr ValidatorTestRequest) error {
	link, params := p.buildURL("problem.saveValidatorTest", url.Values(vtr))
	_, err := p.makeQuery(http.MethodPost, link, params)

	return err
}

func (p *Polygon) SaveCheckerTest(ctr CheckerTestRequest) error {
	link, params := p.buildURL("problem.saveCheckerTest", url.Values(ctr))
	_, err := p.makeQuery(http.MethodPost, link, params)

	return err
}

func (p *Polygon) SaveTestGroup(tgr TestGroupRequest) error {
	link, params := p.buildURL("problem.saveTestGroup", url.Values(tgr))
	_, err := p.makeQuery(http.MethodPost, link, params)

	return err
}

func (p *Polygon) SaveStatementResource(pID int, name, data string) error {
	link, params := p.buildURL("problem.saveStatementResource", url.Values{
		"problemId": {strconv.Itoa(pID)},
		"name":      {name},
		"file":      {data},
	})
	_, err := p.makeQuery(http.MethodPost, link, params)

	return err
}
