package ejudge

import (
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"

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
	url      string
	login    string
	password string

	client *http.Client
}

func NewEjudge(cfg *Config) *Ejudge {
	jar, _ := cookiejar.New(nil)
	return &Ejudge{
		url:      cfg.URL,
		login:    cfg.Login,
		password: cfg.Password,
		client: &http.Client{
			Jar: jar,
		},
	}
}

func (ej *Ejudge) postRequest(method string, params url.Values) (*http.Request, *goquery.Document, error) {
	url, err := url.JoinPath(ej.url, method)
	if err != nil {
		return nil, nil, err
	}
	logrus.WithField("url", url).Debug("post query")
	resp, err := ej.client.PostForm(url, params) //nolint:noctx  // don't need context here
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
		"login":    {ej.login},
		"password": {ej.password},
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

func (ej *Ejudge) CheckContest(sid, cid string) error {
	logrus.WithFields(logrus.Fields{"CID": cid, "SID": sid}).Info("check contest settings, wait please")
	_, doc, err := ej.postRequest("serve-control", url.Values{
		"contest_id": {cid},
		"SID":        {sid},
		"action":     {"262"},
	})
	if err != nil {
		return err
	}
	logrus.WithFields(logrus.Fields{"CID": cid, "SID": sid}).Infof("ejudge answer %q", doc.Find("h2").Text())
	return nil
}
