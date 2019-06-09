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
func NewIssueController() *IssueController{
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

func (ic *IssueController) NewIssueAction(rp core.Roundtrip) {
	//校验session判断是否登陆
	_, err := session(rp)
	if err != nil {
		rp.Redirect("/login")
	} else {
		rp.Display("issue/new")
	}
}

func (ic *IssueController) CreateIssueAction(rp core.Roundtrip) {
	//校验session判断是否登陆
	sess, err := session(rp)
	if err != nil {
		rp.Redirect("/login")
		return
	}

	user, err := sess.User()
	if err != nil {
		msg := fmt.Sprintf("Can't got user... %v", err)
		rp.Redirect(fmt.Sprintf("/err?msg=%s", msg))
		return
	}

	topic := rp.Param("topic")
	_, err = user.CreateIssue(topic)
	if err != nil {
		msg := fmt.Sprintf("Can't got user... %v", err)
		rp.Redirect(fmt.Sprintf("/err?msg=%s", msg))
		return
	}

	rp.Redirect("/index")
}

func (ic *IssueController) ReadIssueAction(rp core.Roundtrip) {
	uuid := rp.Param("id")
	issue, err := model.GetIssueByUUID(uuid)
	if err != nil {
		msg := fmt.Sprintf("Can't got issue... %v", err)
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

	rp.Assign("issue", &issue)
	//issue.Replies()
	//fmt.Println("A-------------------", helper.JsonEncode(issue.Replies()))
	rp.Display("issue/detail")
}

func (ic *IssueController) ReplyIssueAction(rp core.Roundtrip) {
	sess, err := session(rp)
	if err != nil {
		rp.Redirect("/login")
		return
	}

	user, err := sess.User()
	if err != nil {
		msg := fmt.Sprintf("Can't got user... %v", err)
		rp.Redirect(fmt.Sprintf("/err?msg=%s", msg))
		return
	}

	uuid := rp.Param("uuid")
	body := rp.Param("body")

	issue, err := model.GetIssueByUUID(uuid)
	if err != nil {
		msg := fmt.Sprintf("Can't got issue... %v", err)
		rp.Redirect(fmt.Sprintf("/err?msg=%s", msg))
		return
	}

	_, err = user.CreateReply(issue, body)
	if err != nil {
		msg := fmt.Sprintf("Can't CreateReply... %v", err)
		rp.Redirect(fmt.Sprintf("/err?msg=%s", msg))
		return
	}


	rp.Redirect(fmt.Sprintf("/issue/read?id=%v", uuid))
}