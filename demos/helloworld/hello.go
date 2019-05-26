package main

import (
	"fmt"
	"github.com/hq-cml/spider-face"
	"github.com/hq-cml/spider-face/core"
	"time"
)

//TODO 使用不够友好简单
type HelloController struct {
	core.BaseController
}

func (hello *HelloController) HelloAction() {
	hello.Echo("hello world!")
}

func (hello *HelloController) IndexAction() {
	hello.Assign("nowtime", time.Now())
	hello.Assign("title", "welcome to spider~")
	hello.Assign("id", hello.Param("id"))
	hello.Display()
}

func (hello *HelloController) GetRouter() []core.ControllerRouter {
	return []core.ControllerRouter{
		{Method:"GET", Pattern:"/hello/:id", Action: "IndexAction",},
		{Method:"GET", Pattern: "/hello", Action:"HelloAction",},
		{Method:"GET", Pattern: "/index", Action:"IndexAction",},
		{Method:"GET", Pattern: "/index/:id", Action:"IndexAction",},
		{Method:"GET", Pattern: "/index/*", Action:"IndexAction",},    //TODO 这种方式不够科学
	}
}

var controllerMap = map[string]core.SpiderController{
	"hello": &HelloController{},
}

func main() {
	//server config
	sConfig := &core.SpiderConfig{
		BindAddr: ":9529",    //监听地址:端口
	}

	//生成实例
	spd, err := spider.NewSpider(sConfig, controllerMap, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	//Run
	spd.Run()
}
