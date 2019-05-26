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

func (hello *HelloController) GetRouter() map[string]interface{} {
	return map[string]interface{} {
		"/hello/:id": "IndexAction",
		"/hello": "HelloAction",
		"/index": "IndexAction",
		"/index/:id": "IndexAction",
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
