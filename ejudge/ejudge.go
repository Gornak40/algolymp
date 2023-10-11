package ejudge

import (
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"

	"github.com/PuerkitoBio/goquery"
	"github.com/sirupsen/logrus"
)

const BadSID = "0000000000000000"

type Config struct {
	URL      string `json:"url"`
	Login    string `json:"login"`
	Password string `json:"password"`
}

type Ejudge struct {
	cfg    *Config
	client *http.Client
}

func NewEjudge(cfg *Config) *Ejudge {
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

func (ej *Ejudge) Commit(sid string, cid int) error {
	logrus.WithFields(logrus.Fields{"CID": cid, "SID": sid}).Info("lock contest for editing")
	_, _, err := ej.postRequest("serve-control", url.Values{
		"contest_id": {strconv.FormatInt(int64(cid), 10)},
		"SID":        {sid},
		"action":     {"276"},
	})
	if err != nil {
		return err
	}
	logrus.WithFields(logrus.Fields{"CID": cid, "SID": sid}).Info("commit changes")
	_, doc, err := ej.postRequest("serve-control", url.Values{
		"SID":    {sid},
		"action": {"303"},
	})
	if err != nil {
		return err
	}
	status := doc.Find("h2").First().Text()
	logrus.WithFields(logrus.Fields{"CID": cid, "SID": sid}).Infof("ejudge answer %q", status)
	return nil
}

func (ej *Ejudge) CheckContest(sid string, cid int, verbose bool) error {
	logrus.WithFields(logrus.Fields{"CID": cid, "SID": sid}).Info("check contest settings, wait please")
	_, doc, err := ej.postRequest("serve-control", url.Values{
		"contest_id": {strconv.FormatInt(int64(cid), 10)},
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

func (ej *Ejudge) ReloadConfig(sid string, cid int) error {
	req, _, err := ej.postRequest("new-master", url.Values{
		"contest_id": {strconv.FormatInt(int64(cid), 10)},
		"SID":        {sid},
		"action":     {"3"},
	})
	if err != nil {
		return err
	}
	csid := req.URL.Query().Get("SID")
	logrus.WithFields(logrus.Fields{"CID": cid, "CSID": csid, "SID": sid}).Info("success master login")
	_, _, err = ej.postRequest("new-master", url.Values{
		"SID":    {csid},
		"action": {"62"},
	})
	if err != nil {
		return err
	}
	logrus.WithFields(logrus.Fields{"CID": cid, "CSID": csid, "SID": sid}).Info("success reload config")
	return nil
}
