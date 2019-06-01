package main

import (
	"github.com/hq-cml/spider-face/core"
	"github.com/hq-cml/spider-face"
	"fmt"
)

/*
  // Handler
  func hello(c echo.Context) error {
    return c.String(http.StatusOK, "Hello, World!")
  }

  func main() {
    // Echo instance
    e := echo.New()

    // Middleware
    e.Use(middleware.Logger())
    e.Use(middleware.Recover())

    // Routes
    e.GET("/", hello)

    // Start server
    e.Logger.Fatal(e.Start(":1323"))
  }
 */

type HelloController struct {
	core.RuntimeController
}

func (hello *HelloController) HelloAction() {
	hello.Echo("hello world!")
}

func (hello *HelloController) GetRouter() []core.ControllerRouter {
	return []core.ControllerRouter{
		{Method:"GET", Pattern: "/hello", Action:"HelloAction",},
	}
}

//func (hello *HelloController) JsonAction() {
//	if hello.Param("encode") == "yes" {
//		hello.OutputJson(map[string]string {
//			"a":"中文",
//			"b":"yingwen",
//		}, true)
//	} else {
//		hello.OutputJson(map[string]string {
//			"a":"中文",
//			"b":"yingwen",
//		})
//	}
//
//}
//
//func (hello *HelloController) IndexAction() {
//	hello.Assign("nowtime", time.Now())
//	hello.Assign("title", "welcome to spider~")
//	hello.Assign("id", hello.Param("id"))
//	hello.Assign("name", hello.Param("name"))
//	hello.Assign("age", hello.Param("age"))
//	hello.Display()
//}
//
//func (hello *HelloController) PostAction() {
//	hello.Assign("nowtime", time.Now())
//	hello.Assign("title", "welcome to spider~")
//	hello.Assign("id", hello.Param("id"))
//	hello.Assign("name", hello.Param("name"))
//	hello.Assign("age", hello.Param("age"))
//	hello.Display("hello/index")
//}
//
//func (hello *HelloController) GetRouter() []core.ControllerRouter {
//	return []core.ControllerRouter{
//		{Method:"GET", Pattern:"/hello/:id", Action: "IndexAction",},
//		{Method:"GET", Pattern: "/hello", Action:"HelloAction",},
//		{Method:"GET", Pattern: "/index", Action:"IndexAction",},
//		{Method:"GET", Pattern: "/index/:id", Action:"IndexAction",},
//		{Method:"GET", Pattern: "/index/*", Action:"IndexAction",},    //TODO 这种方式不够科学
//		{Method:"POST", Pattern: "/index/post", Action:"PostAction",},
//		{Method:"GET", Pattern: "/json", Action:"JsonAction",},
//	}
//}

var controllers = []core.Controller{
	&HelloController{},
}

func main() {
	//server config
	sConfig := &core.SpiderConfig{
		BindAddr: ":9529",    //监听地址:端口
	}

	//生成实例
	spd, err := spider.NewSpider(sConfig, controllers, nil)
	if err != nil {
		fmt.Println(err)
		return
	}


	//spd.GET("/hello", func() {
	//	hello.Echo("hello world!")
	//})

	//Run
	spd.Run()
}
