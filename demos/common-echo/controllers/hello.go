package controllers

import (
	"github.com/hq-cml/spider-face/core"
	"net/http"
	"fmt"
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
		{Method: http.MethodGet,  Location: "/index",		 Action:"IndexAction",},    //能接收普通参数
		{Method: http.MethodGet,  Location: "/index/:id", 	 Action:"IndexAction",},    //能接收普通参数和路径参数
		{Method: http.MethodPost, Location: "/index/post",   Action:"PostAction",},     //能接收Post参数
		{Method: http.MethodGet,  Location: "/index/" + core.PATH_INFO_IDENTITY,
				Action:"IndexAction",},                                                 //能接收普通参数和pathinfo参数
		{Method: http.MethodGet,  Location: "/json",         Action:"JsonAction",},     //输出Json
	}

	return hc
}

//创建HelloController所拥有的Action
//curl 'http://192.168.110.133:9529/index/aaa'
//curl 'http://192.168.110.133:9529/index/name/hq/age/28'
func (hello *HelloController) IndexAction(rp core.Roundtrip) {
	rp.Echo(fmt.Sprintf("当前时间: %v\n", time.Now()))
	rp.Echo(fmt.Sprintf("参数Id: %v\n", rp.Param("id")))
	rp.Echo(fmt.Sprintf("参数name: %v\n", rp.Param("name")))
	rp.Echo(fmt.Sprintf("参数age: %v\n", rp.Param("age")))
}

//curl -X POST 'http://192.168.110.133:9529/index/post' -d 'id=aaa'
func (hello *HelloController) PostAction(rp core.Roundtrip) {
	rp.Echo(fmt.Sprintf("参数Id: %v\n", rp.Param("id")))
}

//输出结构化的Json
func (hello *HelloController) JsonAction(rp core.Roundtrip) {
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
