package main

import (
	"fmt"
	spider "github.com/hq-cml/spider-face"
	"github.com/hq-cml/spider-face/core"
	"time"
)

type HelloController struct {
	core.BaseController
}

func (hello *HelloController) HelloAction() {
	hello.Echo("hello world!")
}

func (hello *HelloController) IndexAction() {
	hello.Assign("nowtime", time.Now())
	hello.Assign("title", "welcome to spider~")
	//demo.Assign("id", demo.Param("id"))   //TODO
	hello.Display()
}

func (hello *HelloController) GetRouter() map[string]interface{} {
	return map[string]interface{}{
		//"/hello/:id": "IndexAction", //TODO
		"/hello": "HelloAction",
		"/index": "IndexAction",
	}
}

var controllerMap = map[string]core.SpiderController{
	"hello": &HelloController{},
}

func main() {
	//server config
	sConfig := &spider.SpiderConfig{
		DocRoot:         "/tmp/face/www",    //静态文件目录
		BindAddr:        ":9529",            //监听地址:端口
		ViewPath:        "/tmp/face/views",  //模板目录
	}

	spd, err := spider.NewSpider(sConfig)
	if err != nil {
		fmt.Println(err)
		return
	}

	//注册控制器
	spd.RegisterController(controllerMap)

	//Run
	spd.Run()
}
