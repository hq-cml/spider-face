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
		{Method: http.MethodGet, Location: "/",      		 Action:"IndexAction",},
		{Method: http.MethodGet, Location: "/index", 		 Action:"IndexAction",},
		{Method: http.MethodGet, Location: "/issue/new", 	 Action:"NewIssueAction",},
		{Method: http.MethodPost, Location: "/issue/create", Action:"CreateIssueAction",},
	})

	//创建user controller，并绑定路由
	uc := controllers.NewUserAction()
	uc.SetRouteEntries([]core.RouteEntry{
		{Method: http.MethodGet,  Location: "/login",  		   Action:"LoginAction",},
		{Method: http.MethodGet,  Location: "/logout", 		   Action:"LogoutAction",},
		{Method: http.MethodGet,  Location: "/signup", 		   Action:"SignupAction",},
		{Method: http.MethodPost, Location: "/signup_account", Action:"SignupAccountAction",},
		{Method: http.MethodPost, Location: "/authenticate",   Action:"AuthenticateAction",},
	})

	//快捷注册一个通用的错误页面
	spd.GET("/err", controllers.Err)

	//注册controller
	if err := spd.RegisterController([]core.Controller{
		ic, uc,
	}); err != nil {
		fmt.Println(err)
		return
	}

	//Run
	spd.Run()
}