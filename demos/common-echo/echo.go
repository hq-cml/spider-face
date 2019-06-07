package main

/*
 * 利用通用注册方式实现一个api项目
 * 注册一个标准的controller，该controller拥有5个路由规则，每个规则都对应有Action（接口逻辑）
 * 并且，这5个规则拥不同的参数接收方式
 *
 * 如果是一个相对大型的项目，需要考虑组织结构，提供http的接口，通用注册是合适的选择
 */
import (
	"fmt"
	"github.com/hq-cml/spider-face/core"
	"github.com/hq-cml/spider-face"
	"github.com/hq-cml/spider-face/demos/common-echo/controllers"
	"net/http"
)

func main() {
	spd := spider.NewSpider(&core.SpiderConfig{    //生成Spider实例
		BindAddr: ":9529",    					   //监听地址:端口
	}, nil)

	hc := controllers.NewHelloAction()                  //创建需要持有的controller，并绑定路由
	hc.SetRouteEntries([]core.RouteEntry{
		{Method: http.MethodGet,  Location: "/index",		 Action:"IndexAction",},    //能接收普通参数
		{Method: http.MethodGet,  Location: "/index/:id", 	 Action:"IndexAction",},    //能接收普通参数和路径参数
		{Method: http.MethodPost, Location: "/index/post",   Action:"PostAction",},     //能接收Post参数
		{Method: http.MethodGet,  Location: "/index/" + core.PATH_INFO_IDENTITY,
			Action:"IndexAction",},                                                     //能接收普通参数和pathinfo参数
		{Method: http.MethodGet,  Location: "/json",         Action:"JsonAction",},     //输出Json
	})

	if err = spd.RegisterController([]core.Controller{  //注册controller
		hc,
	}); err != nil {
		fmt.Println(err)
		return
	}

	spd.Run()                                           //Run
}
