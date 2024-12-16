package ejudge

import (
	"net/url"
	"strconv"

	"github.com/sirupsen/logrus"
)

func (ej *Ejudge) RegisterUser(csid, login string) error {
	logrus.WithFields(logrus.Fields{"CSID": csid, "login": login}).Info("register user (pending)")
	_, _, err := ej.postRequest(newMaster, url.Values{
		"SID":       {csid},
		"action":    {"20"},
		"add_login": {login},
	})

	return err
}

func (ej *Ejudge) FlipUserVisible(csid string, uid int) error {
	logrus.WithFields(logrus.Fields{"CSID": csid, "uid": uid}).Info("flip user invisible")
	_, _, err := ej.postRequest(newMaster, url.Values{
		"SID":     {csid},
		"action":  {"121"},
		"user_id": {strconv.Itoa(uid)},
	})

	return err
}

func (ej *Ejudge) FlipUserBan(csid string, uid int) error {
	logrus.WithFields(logrus.Fields{"CSID": csid, "uid": uid}).Info("flip user banned")
	_, _, err := ej.postRequest(newMaster, url.Values{
		"SID":     {csid},
		"action":  {"122"},
		"user_id": {strconv.Itoa(uid)},
	})

	return err
}

func (ej *Ejudge) FlipUserLock(csid string, uid int) error {
	logrus.WithFields(logrus.Fields{"CSID": csid, "uid": uid}).Info("flip user locked")
	_, _, err := ej.postRequest(newMaster, url.Values{
		"SID":     {csid},
		"action":  {"123"},
		"user_id": {strconv.Itoa(uid)},
	})

	return err
}

func (ej *Ejudge) FlipUserIncom(csid string, uid int) error {
	logrus.WithFields(logrus.Fields{"CSID": csid, "uid": uid}).Info("flip user incomplete")
	_, _, err := ej.postRequest(newMaster, url.Values{
		"SID":     {csid},
		"action":  {"124"},
		"user_id": {strconv.Itoa(uid)},
	})

	return err
}

func (ej *Ejudge) FlipUserPriv(csid string, uid int) error {
	logrus.WithFields(logrus.Fields{"CSID": csid, "uid": uid}).Info("flip user privileged")
	_, _, err := ej.postRequest(newMaster, url.Values{
		"SID":     {csid},
		"action":  {"289"},
		"user_id": {strconv.Itoa(uid)},
	})

	return err
}
