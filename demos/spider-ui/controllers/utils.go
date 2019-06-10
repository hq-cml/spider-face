package controllers

import "github.com/hq-cml/spider-face/core"

func Err(rp core.Roundtrip) {
	msg := rp.Param("msg")
	rp.Assign("Msg", msg)

	rp.Display("public/error")
}