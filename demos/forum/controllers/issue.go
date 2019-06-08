package controllers

import (
	"github.com/hq-cml/spider-face/core"

	"github.com/hq-cml/spider-face/demos/forum/model"
	"fmt"
)

//创建标准的Controller，必须以Controller后缀结尾
//并实现core.Controller接口
type IssueController struct {
	spdrp core.SpiderRoundtrip
	entries []core.RouteEntry
}
func (ic *IssueController) GetAllRouters() []core.RouteEntry {
	return ic.entries
}
func (ic *IssueController) GetRoundTrip() core.Roundtrip {
	return &ic.spdrp
}

//创建一个controller
func NewIssueAction() *IssueController{
	hc := &IssueController{}
	return hc
}

//设置路由规则
func (ic *IssueController) SetRouteEntries(entries []core.RouteEntry) {
	ic.entries = entries
}

//首页展示
func (ic *IssueController) IndexAction(rp core.Roundtrip) {
	issues, err := model.GetAllIssues()
	if err != nil {
		msg := fmt.Sprintf("Can't got any issues... %v", err)
		rp.Redirect(fmt.Sprintf("/err?msg=%s", msg))
		return
	}

	//校验session判断是否登陆
	_, err = session(rp)
	if err != nil {
		rp.Assign("login", false)
	} else {
		rp.Assign("login", true)
	}

	rp.Assign("issues", issues)
	rp.Display("issue/index")
}
