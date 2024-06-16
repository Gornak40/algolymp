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
	defaultTestset   = "tests"
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
				if err := writer.WriteField(k, v); err != nil {
					return nil, err
				}
			}
		}
		writer.Close()
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
	pairs := []string{}
	for k, vals := range params {
		for _, v := range vals {
			pairs = append(pairs, fmt.Sprintf("%s=%s", k, v))
		}
	}
	sort.Strings(pairs)

	return strings.Join(pairs, "&")
}

func (p *Polygon) buildURL(method string, params url.Values) (string, url.Values) {
	url, _ := url.JoinPath(p.cfg.URL, "api", method)

	params["apiKey"] = []string{p.cfg.APIKey}
	params["time"] = []string{strconv.FormatInt(time.Now().Unix(), 10)}
	sig := fmt.Sprintf("%s/%s?%s#%s", sixSecretSymbols, method, p.skipEscape(params), p.cfg.APISecret)

	b := sha512.Sum512([]byte(sig))
	hsh := hex.EncodeToString(b[:])
	params["apiSig"] = []string{sixSecretSymbols + hsh}

	return url, params
}

func (p *Polygon) BuildPackage(pID int, full, verify bool) error {
	link, params := p.buildURL("problem.buildPackage", url.Values{
		"problemId": []string{strconv.Itoa(pID)},
		"full":      []string{strconv.FormatBool(full)},
		"verify":    []string{strconv.FormatBool(verify)},
	})
	_, err := p.makeQuery(http.MethodPost, link, params)

	return err
}

// Problem Idx (A, B, C) -> Problem.
func (p *Polygon) ContestProblems(pID int) (map[string]ProblemAnswer, error) {
	link, params := p.buildURL("contest.problems", url.Values{
		"contestId": []string{strconv.Itoa(pID)},
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
		"problemId":    []string{strconv.Itoa(pID)},
		"minorChanges": []string{strconv.FormatBool(minor)},
		"message":      []string{message},
	})
	_, err := p.makeQuery(http.MethodPost, link, params)

	return err
}

func (p *Polygon) UpdateWorkingCopy(pid int) error {
	link, params := p.buildURL("problem.updateWorkingCopy", url.Values{
		"problemId": []string{strconv.Itoa(pid)},
	})
	_, err := p.makeQuery(http.MethodPost, link, params)

	return err
}

func (p *Polygon) GetPackages(pID int) ([]PackageAnswer, error) {
	link, params := p.buildURL("problem.packages", url.Values{
		"problemId": []string{strconv.Itoa(pID)},
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
		"problemId": []string{strconv.Itoa(pID)},
		"testset":   []string{defaultTestset},
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
		"id": []string{strconv.Itoa(pID)},
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
		"problemId": []string{strconv.Itoa(pID)},
		"testset":   []string{defaultTestset},
		"noInputs":  []string{"true"},
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
		"problemId": []string{strconv.Itoa(pID)},
		"packageId": []string{strconv.Itoa(packID)},
		"type":      []string{packType},
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
		"problemId": []string{strconv.Itoa(pID)},
		"testset":   []string{defaultTestset},
		"enable":    []string{"true"},
	})
	_, err := p.makeQuery(http.MethodPost, link, params)

	return err
}

func (p *Polygon) EnablePoints(pID int) error {
	link, params := p.buildURL("problem.enablePoints", url.Values{
		"problemId": []string{strconv.Itoa(pID)},
		"enable":    []string{"true"},
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
		"problemId": []string{strconv.Itoa(pID)},
		"tags":      []string{tags},
	})
	_, err := p.makeQuery(http.MethodPost, link, params)

	return err
}

func (p *Polygon) SetValidator(pID int, validator string) error {
	link, params := p.buildURL("problem.setValidator", url.Values{
		"problemId": []string{strconv.Itoa(pID)},
		"validator": []string{validator},
	})
	_, err := p.makeQuery(http.MethodPost, link, params)

	return err
}

func (p *Polygon) SetChecker(pID int, checker string) error {
	link, params := p.buildURL("problem.setChecker", url.Values{
		"problemId": []string{strconv.Itoa(pID)},
		"checker":   []string{checker},
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
		"problemId":  []string{strconv.Itoa(pID)},
		"interactor": []string{interactor},
	})
	_, err := p.makeQuery(http.MethodPost, link, params)

	return err
}

func (p *Polygon) SaveSolution(sr SolutionRequest) error {
	link, params := p.buildURL("problem.saveSolution", url.Values(sr))
	_, err := p.makeQuery(http.MethodPost, link, params)

	return err
}
