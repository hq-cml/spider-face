package main

import (
	"github.com/hq-cml/spider-face"
	"github.com/hq-cml/spider-face/core"
	"net/http"
	"github.com/hq-cml/spider-face/demos/forum/controllers"
	"fmt"
)

func main() {
	spd := spider.NewSpider(&core.SpiderConfig{         //生成Spider实例
		BindAddr: ":9529",    						    //监听地址:端口
	}, nil)

	//创建issue controller，并绑定路由
	ic := controllers.NewIssueAction()
	ic.SetRouteEntries([]core.RouteEntry{
		{Method: http.MethodGet, Location: "/",      Action:"IndexAction",},
		{Method: http.MethodGet, Location: "/index", Action:"IndexAction",},
	})

	//创建user controller，并绑定路由



	if err := spd.RegisterController([]core.Controller{  //注册controller
		ic,
	}); err != nil {
		fmt.Println(err)
		return
	}

	//Run
	spd.Run()
}