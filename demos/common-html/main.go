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
	"github.com/hq-cml/spider-face/demos/common-html/controllers"
	"net/http"
)

func main() {
	spd := spider.NewSpider(&core.SpiderConfig{         //生成Spider实例
		BindAddr: ":9529",    						    //监听地址:端口
	}, nil)


	hc := controllers.NewHelloAction()                  //创建需要持有的controller，并绑定路由
	hc.SetRouteEntries([]core.RouteEntry{
		{Method: http.MethodGet, Location: "/index",     Action:"IndexAction",},
		{Method: http.MethodGet, Location: "/index/:id", Action:"IndexAction",},
	})

	if err := spd.RegisterController([]core.Controller{  //注册controller
		hc,
	}); err != nil {
		fmt.Println(err)
		return
	}

	spd.Run()                                            //Run
}
