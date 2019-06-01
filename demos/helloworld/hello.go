package main

import (
	"github.com/hq-cml/spider-face/core"
	"github.com/hq-cml/spider-face"
	"fmt"
	"time"
)

type HelloController struct {
	core.RuntimeController
}

func (hello *HelloController) HelloAction() {
	hello.Echo("hello world!")
}

//func (hello *HelloController) GetRouter() []core.ControllerRouter {
//	return []core.ControllerRouter{
//		{Method:"GET", Location: "/hello", Action:"HelloAction",},
//	}
//}

func (hello *HelloController) JsonAction() {
	if hello.Param("encode") == "yes" {
		hello.OutputJson(map[string]string {
			"a":"中文",
			"b":"yingwen",
		}, true)
	} else {
		hello.OutputJson(map[string]string {
			"a":"中文",
			"b":"yingwen",
		})
	}
}

func (hello *HelloController) IndexAction() {
	hello.Assign("nowtime", time.Now())
	hello.Assign("title", "welcome to spider~")
	hello.Assign("id", hello.Param("id"))
	hello.Assign("name", hello.Param("name"))
	hello.Assign("age", hello.Param("age"))
	hello.Display()
}

func (hello *HelloController) PostAction() {
	hello.Assign("nowtime", time.Now())
	hello.Assign("title", "welcome to spider~")
	hello.Assign("id", hello.Param("id"))
	hello.Assign("name", hello.Param("name"))
	hello.Assign("age", hello.Param("age"))
	hello.Display("hello/index")
}

func (hello *HelloController) GetRouter() []core.ControllerRouter {
	return []core.ControllerRouter{
		{Method:"GET", Location:"/hello/:id", Action: "IndexAction",},
		{Method:"GET", Location: "/hello", Action:"HelloAction",},
		{Method:"GET", Location: "/index", Action:"IndexAction",},
		{Method:"GET", Location: "/index/:id", Action:"IndexAction",},
		{Method:"GET", Location: "/index/*", Action:"IndexAction",},    //TODO 这种方式不够科学
		{Method:"POST", Location: "/index/post", Action:"PostAction",},
		{Method:"GET", Location: "/json", Action:"JsonAction",},
	}
}

var controllers = []core.Controller{
	&HelloController{},
}

func main() {
	//server config
	sConfig := &core.SpiderConfig{
		BindAddr: ":9529",    //监听地址:端口
	}

	//生成实例
	spd, err := spider.NewSpider(sConfig, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	//注册controller
	err = spd.RegisterController(controllers)
	if err != nil {
		fmt.Println(err)
		return
	}

	//spd.GET("/hello", func(c core.Controller) {
	//	fmt.Println("M--------------", c.Param("id"))
	//	c.Echo("Hello")
	//})

	//快捷注册
	spd.GET("/fast/get", func(c core.Controller) {
		c.Echo("参数Id是：" + c.Param("id"))
	})

	//spd.GET("/fast/display", func(c core.Controller) {
	//	fmt.Println("F--------------", c.Param("id"))
	//	c.Echo("参数Id是：" + c.Param("id"))
	//	c.Assign("nowtime", time.Now())
	//	c.Assign("title", "welcome to spider~")
	//	c.Assign("id", c.Param("id"))
	//	c.Assign("name", c.Param("name"))
	//	c.Assign("age", c.Param("age"))
	//
	//	c.Display("hello/index")
	//})

	spd.POST("/fast/post", func(c core.Controller) {
		c.Echo("参数Id是：" + c.Param("id"))
	})

	//Run
	spd.Run()
}
