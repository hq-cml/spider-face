package main

import (
	"github.com/hq-cml/spider-face"
	"github.com/hq-cml/spider-face/core"
	"net/http"
	"github.com/hq-cml/spider-face/demos/spider-ui/controllers"
	"fmt"
)

func main() {
	spd := spider.NewSpider(&core.SpiderConfig{         //生成Spider实例
		BindAddr: ":9530",    						    //监听地址:端口
	}, nil)

	//创建issue controller，并绑定路由
	wc := controllers.NewWeiboController()
	wc.SetRouteEntries([]core.RouteEntry{
		{Method: http.MethodGet,  Location: "/",      		 Action:"IndexAction",},
		{Method: http.MethodGet,  Location: "/index", 		 Action:"IndexAction",},
		{Method: http.MethodPost, Location: "/_search", 	 Action:"SearchAction",},
		{Method: http.MethodGet,  Location: "/_detail",      Action:"DetailAction",},
	})

	controllers.SpiderEngineAddr = "192.168.110.133:9528"
	controllers.SpiderEngineDb = "sp_db"
	controllers.SpiderEngineTable = "weibo"

	//注册controller
	if err := spd.RegisterController([]core.Controller{
		wc,
	}); err != nil {
		fmt.Println(err)
		return
	}

	//快捷注册一个通用的错误页面
	spd.GET("/err", controllers.Err)

	//Run
	spd.Run()
}
