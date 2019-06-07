package controllers

import (
	"github.com/hq-cml/spider-face/core"
	"net/http"
	"time"
)

//创建标准的Controller，必须以Controller后缀结尾
//并实现core.Controller接口
type HelloController struct {
	spdrp core.SpiderRoundtrip
	entries []core.RouteEntry
}
func (hello *HelloController) GetAllRouters() []core.RouteEntry {
	return hello.entries
}
func (hello *HelloController) GetRoundTrip() core.Roundtrip {
	return &hello.spdrp
}

//创建一个controller，并绑定路由
func NewHelloAction() *HelloController{
	hc := &HelloController{}
	hc.entries = []core.RouteEntry{
		{Method: http.MethodGet, Location: "/index",     Action:"IndexAction",},
		{Method: http.MethodGet, Location: "/index/:id", Action:"IndexAction",},
	}

	return hc
}

//浏览器 'http://192.168.110.133:9529/index/aaa'
//浏览器 'http://192.168.110.133:9529/index/name/hq/age/28'
func (hello *HelloController) IndexAction(rp core.Roundtrip) {
	rp.Assign("nowtime", time.Now())
	rp.Assign("title", "welcome to spider~")
	rp.Assign("id", rp.Param("id"))
	rp.Assign("name", rp.Param("name"))
	rp.Assign("age", rp.Param("age"))
	rp.Display()
}
