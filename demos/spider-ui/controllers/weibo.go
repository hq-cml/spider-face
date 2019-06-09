package controllers

import (
	"github.com/hq-cml/spider-face/core"
)

//创建标准的Controller，必须以Controller后缀结尾
//并实现core.Controller接口
type WeiboController struct {
	spdrp core.SpiderRoundtrip
	entries []core.RouteEntry
}
func (ic *WeiboController) GetAllRouters() []core.RouteEntry {
	return ic.entries
}
func (ic *WeiboController) GetRoundTrip() core.Roundtrip {
	return &ic.spdrp
}

//创建一个controller
func NewWeiboController() *WeiboController {
	hc := &WeiboController{}
	return hc
}

//设置路由规则
func (ic *WeiboController) SetRouteEntries(entries []core.RouteEntry) {
	ic.entries = entries
}

//首页展示
func (ic *WeiboController) IndexAction(rp core.Roundtrip) {
	rp.Display("weibo/index")
}
