package controllers

import (
	"github.com/hq-cml/spider-face/core"
	"net/http"
	"strings"
	"io/ioutil"
	"fmt"
	"github.com/hq-cml/spider-face/utils/helper"
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

	//创建自定义的client进行搜索
	client := &http.Client{}
	query := map[string]string{
		"database" : SpiderEngineDb,
		"table" : SpiderEngineTable,
		"value" : keyword,
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
	rp.Echo(string(body))
}