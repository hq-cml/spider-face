package main

import (
	"github.com/hq-cml/spider-face/core"
	"github.com/hq-cml/spider-face"
	"fmt"
	"time"
)

type HelloController struct {
	spdrp core.SpiderRoundtrip
}

func (hello *HelloController) HelloAction(rp core.Roundtrip) {
	rp.Echo("hello world!")
}

func (hello *HelloController) GetAllRouters() []core.ControllerRouter { //TODO 这种方式不够友好
	return []core.ControllerRouter{
		{Method:"GET", Location:"/hello/:id", Action: "IndexAction",},
		{Method:"GET", Location: "/hello", Action:"HelloAction",},
		{Method:"GET", Location: "/index", Action:"IndexAction",},
		{Method:"GET", Location: "/index/:id", Action:"IndexAction",},
		{Method:"GET", Location: "/index/*", Action:"IndexAction",},
		{Method:"POST", Location: "/index/post", Action:"PostAction",},
		{Method:"GET", Location: "/json", Action:"JsonAction",},
	}
}

func (hello *HelloController) GetRoundTrip() core.Roundtrip {
	return &hello.spdrp
}

func (hello *HelloController) JsonAction(rp core.Roundtrip) {
	if rp.Param("encode") == "yes" {
		rp.OutputJson(map[string]string {
			"a":"中文",
			"b":"yingwen",
		}, true)
	} else {
		rp.OutputJson(map[string]string {
			"a":"中文",
			"b":"yingwen",
		})
	}
}

func (hello *HelloController) IndexAction(rp core.Roundtrip) {
	rp.Assign("nowtime", time.Now())
	rp.Assign("title", "welcome to spider~")
	rp.Assign("id", rp.Param("id"))
	rp.Assign("name", rp.Param("name"))
	rp.Assign("age", rp.Param("age"))
	rp.Display()
}

func (hello *HelloController) PostAction(rp core.Roundtrip) {
	rp.Assign("nowtime", time.Now())
	rp.Assign("title", "welcome to spider~")
	rp.Assign("id", rp.Param("id"))
	rp.Assign("name", rp.Param("name"))
	rp.Assign("age", rp.Param("age"))
	rp.Display("hello/index")
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
	err = spd.RegisterController([]core.Controller{
		&HelloController{},
	})
	if err != nil {
		fmt.Println(err)
		return
	}


	//快捷注册
	spd.GET("/fast/get", func(rp core.Roundtrip) {
		rp.Echo("参数Id是：" + rp.Param("id"))
	})

	spd.GET("/fast/display", func(rp core.Roundtrip) {
		rp.Assign("nowtime", time.Now())
		rp.Assign("title", "welcome to spider~")
		rp.Assign("id", rp.Param("id"))
		rp.Assign("name", rp.Param("name"))
		rp.Assign("age", rp.Param("age"))

		rp.Display("hello/index")
	})

	spd.POST("/fast/post", func(rp core.Roundtrip) {
		rp.Echo("参数Id是：" + rp.Param("id"))
	})

	//Run
	spd.Run()
}
