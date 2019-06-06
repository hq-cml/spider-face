package main
/*
 * 利用通用注册方式实现一个web项目
 * 注册一个标准的controller，该controller拥有3个路由规则，每个规则都对应有Action（接口逻辑）
 * 并且，这3个规则拥不同的参数接收方式
 *
 * 如果是一个相对大型的项目，需要考虑组织结构，提供web功能，通用注册是合适的选择
 */

import (
	"github.com/hq-cml/spider-face/core"
	"github.com/hq-cml/spider-face"
	"fmt"
	"time"
	"net/http"
)

//标准的Controller，必须以Controller后缀结尾
type HelloController struct {
	spdrp core.SpiderRoundtrip
}

//实现core.Controller接口
func (hello *HelloController) GetAllRouters() []core.ControllerRouter { //TODO 这种方式不够友好
	return []core.ControllerRouter{
		{Method: http.MethodGet, Location: "/index", Action:"IndexAction",},
		{Method: http.MethodGet, Location: "/index/:id", Action:"IndexAction",},
		{Method: http.MethodGet, Location: "/index/*", Action:"IndexAction",},
	}
}
func (hello *HelloController) GetRoundTrip() core.Roundtrip {
	return &hello.spdrp
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

func main() {
	spd, err := spider.NewSpider(&core.SpiderConfig{    //生成Spider实例
		BindAddr: ":9529",    						    //监听地址:端口
	}, nil)
	if err != nil {
		return
	}

	if err = spd.RegisterController([]core.Controller{  //注册controller
		&HelloController{},
	}); err != nil {
		fmt.Println(err)
		return
	}

	spd.Run()                                           //Run
}
