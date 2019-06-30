package controllers

import (
	"github.com/hq-cml/spider-face/core"
	"net/http"
	"strings"
	"io/ioutil"
	"fmt"
	"github.com/hq-cml/spider-face/utils/helper"
	"encoding/json"
	"time"
	"strconv"
)

var SpiderEngineAddr string
var SpiderEngineDb string
var SpiderEngineTable string

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

//搜索接口
func (ic *WeiboController) SearchAction(rp core.Roundtrip) {
	keyword := rp.Param("keyword")
	p := rp.Param("page")
	page := 1
	size := 10
	var err error
	if p != "" {
		page, err = strconv.Atoi(p)
		if err != nil {
			msg := fmt.Sprintf("Search Error... %v", err)
			rp.Redirect(fmt.Sprintf("/err?msg=%s", msg))
			return
		}
	}

	//创建自定义的client进行搜索
	client := &http.Client{}
	query := map[string]interface{}{
		"database" : SpiderEngineDb,
		"table" : SpiderEngineTable,
		"value" : keyword,
		"offset": size * (page-1),
		"size": size,
	}
	req, _ := http.NewRequest("GET",
		fmt.Sprintf("http://%s/_search", SpiderEngineAddr), strings.NewReader(helper.JsonEncode(query)))
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		msg := fmt.Sprintf("Search Error... %v", err)
		rp.Redirect(fmt.Sprintf("/err?msg=%s", msg))
		return
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	r := Result{}
	err = json.Unmarshal(body, &r)
	if err != nil {
		msg := fmt.Sprintf("Json decode Error... %v", err)
		rp.Redirect(fmt.Sprintf("/err?msg=%s", msg))
		return
	}

	list := []Detail{}
	for _, d := range r.Data.Docs {
		if len([]rune(d.Detail.Content)) > 20 {
			d.Detail.Summary = string([]rune(d.Detail.Content)[0: 20]) + "。。。"
		} else {
			d.Detail.Summary = d.Detail.Content
		}
		d.Detail.CreatedAt = time.Unix(d.Detail.Date, 0)
		list = append(list, d.Detail)
	}

	rp.Assign("list", list)

	rp.Display("weibo/list")
}

type Result struct {
	Code int			`json:"code"`
	Msg  string			`json:"msg"`
	Data DocList	    `json:"data"`
}
type DocList struct {
	Total int `json:"total"`
	Docs  []DocInfo `json:"docs"`
}
type DocInfo struct {
	Key    string
	Detail Detail
}

type Detail struct {
	Date 		int64  `json:"date"`
	ReadCnt 	int64  `json:"read_cnt"`
	User 		string `json:"user_name"`
	Content 	string `json:"weibo_content"`
	Id 			string `json:"weibo_id"`
	Summary     string
	CreatedAt   time.Time
}

//详情接口
func (ic *WeiboController) DetailAction(rp core.Roundtrip) {
	id := rp.Param("id")

	//创建自定义的client进行搜索
	client := &http.Client{}
	req, _ := http.NewRequest("GET",
		fmt.Sprintf("http://%s/%s/%s/%s", SpiderEngineAddr, SpiderEngineDb, SpiderEngineTable, id), nil)

	resp, err := client.Do(req)
	if err != nil {
		msg := fmt.Sprintf("Search Error... %v", err)
		rp.Redirect(fmt.Sprintf("/err?msg=%s", msg))
		return
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	r := DetailResult{}
	err = json.Unmarshal(body, &r)
	if err != nil {
		msg := fmt.Sprintf("Json decode Error... %v", err)
		rp.Redirect(fmt.Sprintf("/err?msg=%s", msg))
		return
	}

	r.Data.Detail.CreatedAt = time.Unix(r.Data.Detail.Date, 0)
	rp.Assign("detail", r.Data.Detail)

	rp.Display("weibo/detail")

}

type DetailResult struct {
	Code int			`json:"code"`
	Msg  string			`json:"msg"`
	Data DocInfo	    `json:"data"`
}