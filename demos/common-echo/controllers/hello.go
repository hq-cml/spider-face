package controllers

import (
	"github.com/hq-cml/spider-face/core"
	"fmt"
	"time"
)

//创建标准的Controller，必须以Controller后缀结尾
//并实现core.Controller接口
type HelloController struct {
	spdrp core.SpiderRoundtrip
	entries []core.RouteEntry
}
func (hc *HelloController) GetAllRouters() []core.RouteEntry {
	return hc.entries
}
func (hc *HelloController) GetRoundTrip() core.Roundtrip {
	return &hc.spdrp
}

//创建一个controller，并绑定路由
func NewHelloAction() *HelloController{
	hc := &HelloController{}
	return hc
}

//设置路由规则
func (hc *HelloController) SetRouteEntries(entries []core.RouteEntry) {
	hc.entries = entries
}

//创建HelloController所拥有的Action
//curl 'http://192.168.110.133:9529/index/aaa'
//curl 'http://192.168.110.133:9529/index/name/hq/age/28'
func (hc *HelloController) IndexAction(rp core.Roundtrip) {
	rp.Echo(fmt.Sprintf("当前时间: %v\n", time.Now()))
	rp.Echo(fmt.Sprintf("参数Id: %v\n", rp.Param("id")))
	rp.Echo(fmt.Sprintf("参数name: %v\n", rp.Param("name")))
	rp.Echo(fmt.Sprintf("参数age: %v\n", rp.Param("age")))
}

//curl -X POST 'http://192.168.110.133:9529/index/post' -d 'id=aaa'
func (hc *HelloController) PostAction(rp core.Roundtrip) {
	rp.Echo(fmt.Sprintf("参数Id: %v\n", rp.Param("id")))
}

//输出结构化的Json
func (hc *HelloController) JsonAction(rp core.Roundtrip) {
	if rp.Param("encode") == "yes" {
		rp.OutputJson(map[string]string {
			"a":"中文",
			"b":"yingwen",
		}, true)     //utf8编码中文
	} else {
		rp.OutputJson(map[string]string {
			"a":"中文",
			"b":"yingwen",
		})
	}
}
