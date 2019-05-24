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
	hello.Assign("id", hello.Param("id"))
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

	//生成实例
	spd, err := spider.NewSpider(sConfig, controllerMap)
	if err != nil {
		fmt.Println(err)
		return
	}

	//Run
	spd.Run()
}
