package controllers

import (
	"errors"
	"github.com/hq-cml/spider-face/core"
	"github.com/hq-cml/spider-face/demos/forum/model"
)

// Checks if the user is logged in and has a session, if not err is not nil
func session(rp core.Roundtrip) (sess model.Session, err error) {
	cookie := rp.GetCookie("_cookie")
	if err == nil {
		sess = model.Session{Uuid: cookie}
		if ok, _ := sess.Check(); !ok {
			err = errors.New("Invalid session")
		}
	}
	return
}