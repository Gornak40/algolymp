package ejudge

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"strconv"

	"github.com/PuerkitoBio/goquery"
	"github.com/sirupsen/logrus"
)

const BadSID = "0000000000000000"

var (
	ErrParseMasterSID = fmt.Errorf("can't parse master SID")
)

type Config struct {
	URL       string `json:"url"`
	Login     string `json:"login"`
	Password  string `json:"password"`
	JudgesDir string `json:"judgesDir"`
}

type Ejudge struct {
	cfg    *Config
	client *http.Client
}

func NewEjudge(cfg *Config) *Ejudge {
	logrus.WithField("url", cfg.URL).Info("init ejudge engine")
	jar, _ := cookiejar.New(nil)

	return &Ejudge{
		cfg: cfg,
		client: &http.Client{
			Jar: jar,
		},
	}
}

func (ej *Ejudge) postRequest(method string, params url.Values) (*http.Request, *goquery.Document, error) {
	url, err := url.JoinPath(ej.cfg.URL, method)
	if err != nil {
		return nil, nil, err
	}
	logrus.WithField("url", url).Debug("post query")
	resp, err := ej.client.PostForm(url, params) //nolint:noctx  // don't need context here.
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, nil, fmt.Errorf("bad status code %d", resp.StatusCode)
	}
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, nil, err
	}

	return resp.Request, doc, nil
}

func (ej *Ejudge) Login() (string, error) {
	req, _, err := ej.postRequest("serve-control", url.Values{
		"login":    {ej.cfg.Login},
		"password": {ej.cfg.Password},
	})
	if err != nil {
		return BadSID, err
	}
	sid := req.URL.Query().Get("SID")
	if sid == "" {
		return BadSID, err
	}
	logrus.WithField("SID", sid).Info("success login")

	return sid, nil
}

// not sure it's working.
func (ej *Ejudge) Logout(sid string) error {
	_, _, err := ej.postRequest("serve-control", url.Values{
		"SID":    {sid},
		"action": {"55"},
	})
	if err != nil {
		return err
	}
	logrus.WithField("SID", sid).Info("success logout")

	return nil
}

func (ej *Ejudge) Lock(sid string, cid int) error {
	logrus.WithFields(logrus.Fields{"CID": cid, "SID": sid}).Info("lock contest for editing")
	_, _, err := ej.postRequest("serve-control", url.Values{
		"contest_id": {strconv.Itoa(cid)},
		"SID":        {sid},
		"action":     {"276"},
	})

	return err
}

func (ej *Ejudge) Commit(sid string) error {
	logrus.WithFields(logrus.Fields{"SID": sid}).Info("commit changes")
	_, doc, err := ej.postRequest("serve-control", url.Values{
		"SID":    {sid},
		"action": {"303"},
	})
	if err != nil {
		return err
	}
	status := doc.Find("h2").First().Text()
	logrus.WithFields(logrus.Fields{"SID": sid}).Infof("ejudge answer %q", status)

	return nil
}

func (ej *Ejudge) CheckContest(sid string, cid int, verbose bool) error {
	logrus.WithFields(logrus.Fields{"CID": cid, "SID": sid}).Info("check contest settings, wait please")
	_, doc, err := ej.postRequest("serve-control", url.Values{
		"contest_id": {strconv.Itoa(cid)},
		"SID":        {sid},
		"action":     {"262"},
	})
	if err != nil {
		return err
	}
	if verbose {
		logrus.Info(doc.Find("font").Text())
	}
	status := doc.Find("h2").First().Text()
	logrus.WithFields(logrus.Fields{"CID": cid, "SID": sid}).Infof("ejudge answer %q", status)

	return nil
}

func (ej *Ejudge) MasterLogin(sid string, cid int) (string, error) {
	req, _, err := ej.postRequest("new-master", url.Values{
		"contest_id": {strconv.Itoa(cid)},
		"SID":        {sid},
		"action":     {"3"},
	})
	if err != nil {
		return "", err
	}
	csid := req.URL.Query().Get("SID")
	if csid == "" {
		return "", ErrParseMasterSID
	}
	logrus.WithFields(logrus.Fields{"CID": cid, "CSID": csid, "SID": sid}).Info("success master login")

	return csid, nil
}

func (ej *Ejudge) FilterRuns(csid string, filter string, count int) ([]int, error) {
	_, doc, err := ej.postRequest("new-master", url.Values{
		"SID":              {csid},
		"filter_view":      {"1"},
		"filter_expr":      {filter},
		"filter_first_run": {"-1"},
		"filter_last_run":  {strconv.Itoa(-count)},
	})
	if err != nil {
		return nil, err
	}
	ejErr := doc.Find("#container > pre")
	if ejErr.Text() != "" {
		return nil, errors.New(ejErr.Text())
	}
	res := doc.Find("#container > table:nth-child(18) > tbody > tr > td:nth-child(1)")
	digits := regexp.MustCompile("[^0-9]+")
	runsStr := res.Map(func(_ int, s *goquery.Selection) string {
		return digits.ReplaceAllString(s.Text(), "")
	})
	runs := make([]int, 0, len(runsStr))
	for _, s := range runsStr {
		run, err := strconv.Atoi(s)
		if err != nil {
			return nil, err
		}
		runs = append(runs, run)
	}

	return runs, nil
}

func (ej *Ejudge) ReloadConfig(csid string) error {
	_, _, err := ej.postRequest("new-master", url.Values{
		"SID":    {csid},
		"action": {"62"},
	})
	if err != nil {
		return err
	}
	logrus.WithFields(logrus.Fields{"CSID": csid}).Info("success reload config")

	return nil
}

func (ej *Ejudge) CreateContest(sid string, cid int, tid int) error {
	logrus.WithFields(logrus.Fields{"CID": cid, "TID": tid, "SID": sid}).Info("create contest")
	_, doc, err := ej.postRequest("serve-control", url.Values{
		"contest_id": {strconv.Itoa(cid)},
		"SID":        {sid},
		"num_mode":   {"1"},
		"action":     {"259"},
		"templ_mode": {"1"},
		"templ_id":   {strconv.Itoa(tid)},
	})
	if err != nil {
		return err
	}
	status := doc.Find("h2").First().Text()
	if status == "" {
		status = "OK"
	}
	logrus.WithFields(logrus.Fields{"CID": cid, "TID": tid, "SID": sid}).Infof("ejudge answer %q", status)

	return nil
}

func (ej *Ejudge) MakeInvisible(sid string, cid int) error {
	logrus.WithFields(logrus.Fields{"CID": cid, "SID": sid}).Info("make invisible")
	_, _, err := ej.postRequest("serve-control", url.Values{
		"contest_id": {strconv.Itoa(cid)},
		"SID":        {sid},
		"action":     {"6"},
	})

	return err
}

func (ej *Ejudge) MakeVisible(sid string, cid int) error {
	logrus.WithFields(logrus.Fields{"CID": cid, "SID": sid}).Info("make visible")
	_, _, err := ej.postRequest("serve-control", url.Values{
		"contest_id": {strconv.Itoa(cid)},
		"SID":        {sid},
		"action":     {"7"},
	})

	return err
}
