package main

/*
 * 利用通用注册方式实现一个api项目
 * 注册一个标准的controller，该controller拥有5个路由规则，每个规则都对应有Action（接口逻辑）
 * 并且，这5个规则拥不同的参数接收方式
 *
 * 如果是一个相对大型的项目，需要考虑组织结构，提供http的接口，通用注册是合适的选择
 */
import (
	"github.com/hq-cml/spider-face/core"
	"github.com/hq-cml/spider-face"
	"fmt"
	"time"
	"net/http"
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
		{Method: http.MethodGet, Location: "/index", Action:"IndexAction",},      //能接收普通参数
		{Method: http.MethodGet, Location: "/index/:id", Action:"IndexAction",},  //能接收普通参数和路径参数
		{Method: http.MethodGet, Location: "/index/*", Action:"IndexAction",},    //能接收普通参数和pathinfo参数
		{Method: http.MethodPost, Location: "/index/post", Action:"PostAction",}, //能接Post参数
		{Method: http.MethodGet, Location: "/json", Action:"JsonAction",},
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

func main() {
	spd, err := spider.NewSpider(&core.SpiderConfig{    //生成Spider实例
		BindAddr: ":9529",    						    //监听地址:端口
	}, nil)
	if err != nil {
		return
	}

	hc := NewHelloAction()                              //创建需要持有的controller，并绑定路由

	if err = spd.RegisterController([]core.Controller{  //注册controller
		hc,
	}); err != nil {
		fmt.Println(err)
		return
	}

	spd.Run()                                           //Run
}
