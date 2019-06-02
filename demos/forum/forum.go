package main

/*
import (
	"fmt"
	"time"
	"github.com/hq-cml/spider-face"
	"github.com/hq-cml/spider-face/core"
)

//TODO 使用不够友好简单
type HelloController struct {
	core.SpiderRoundtrip
}

func (hello *HelloController) HelloAction() {
	hello.Echo("hello world!")
}

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

func (hello *HelloController) GetAllRouters() []core.ControllerRouter {
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

var controllerMap = map[string]core.Controller{
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
*/