package ejudge

import (
	"bytes"
	"encoding/csv"
	"io"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/sirupsen/logrus"
)

const (
	defBufSize = 1024
)

func (ej *Ejudge) DumpUsers(csid string) (io.Reader, error) {
	logrus.WithFields(logrus.Fields{
		"CSID": csid,
	}).Info("dump contest users")
	_, doc, err := ej.postRequest(newMaster, url.Values{
		"SID":    {csid},
		"action": {"132"},
	})
	if err != nil {
		return nil, err
	}

	return strings.NewReader(doc.Text()), nil // TODO: fix trimspace
}

func (ej *Ejudge) DumpRuns(csid string) (io.Reader, error) {
	logrus.WithFields(logrus.Fields{
		"CSID": csid,
	}).Info("dump contest runs")
	_, doc, err := ej.postRequest(newMaster, url.Values{
		"SID":    {csid},
		"action": {"152"},
	})
	if err != nil {
		return nil, err
	}

	return strings.NewReader(doc.Text()), nil
}

func (ej *Ejudge) DumpStandings(csid string) (io.Reader, error) {
	logrus.WithFields(logrus.Fields{
		"CSID": csid,
	}).Info("dump contest standings")
	_, doc, err := ej.postRequest(newMaster, url.Values{
		"SID":    {csid},
		"action": {"94"},
	})
	if err != nil {
		return nil, err
	}
	th := doc.Find("table.standings > tbody > tr")

	return walkTable(th)
}

func (ej *Ejudge) DumpProbStats(csid string) (io.Reader, error) {
	logrus.WithFields(logrus.Fields{
		"CSID": csid,
	}).Info("dump problem stats")
	_, doc, err := ej.postRequest(newMaster, url.Values{
		"SID":    {csid},
		"action": {"309"},
	})
	if err != nil {
		return nil, err
	}
	th := doc.Find("table.b1 > tbody > tr")

	return walkTable(th)
}

func (ej *Ejudge) DumpRegPasswords(csid string) (io.Reader, error) {
	logrus.WithFields(logrus.Fields{
		"CSID": csid,
	}).Info("dump registration passwords")
	_, doc, err := ej.postRequest(newMaster, url.Values{
		"SID":    {csid},
		"action": {"120"},
	})
	if err != nil {
		return nil, err
	}
	th := doc.Find("table.b1 > tbody > tr")

	return walkTable(th)
}

func (ej *Ejudge) DumpIPs(csid string) (io.Reader, error) {
	logrus.WithFields(logrus.Fields{
		"CSID": csid,
	}).Info("dump user IPs")
	_, doc, err := ej.postRequest(newMaster, url.Values{
		"SID":    {csid},
		"action": {"235"},
	})
	if err != nil {
		return nil, err
	}
	th := doc.Find("table.b1 > tbody > tr")

	return walkTable(th)
}

func walkTable(table *goquery.Selection) (io.Reader, error) {
	bf := bytes.NewBuffer(make([]byte, 0, defBufSize))
	w := csv.NewWriter(bf)
	w.Comma = ';'

	var err error
	table.EachWithBreak(func(_ int, row *goquery.Selection) bool {
		cols := row.Find("th, td")
		rec := make([]string, 0, cols.Length())
		cols.Each(func(_ int, s *goquery.Selection) {
			rec = append(rec, s.Text())
		})
		if err = w.Write(rec); err != nil {
			return false
		}

		return true
	})
	w.Flush()

	return strings.NewReader(bf.String()), err
}
