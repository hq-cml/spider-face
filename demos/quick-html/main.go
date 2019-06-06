package main

/*
 * 利用Spider的快捷注册功能，快速开始一个项目
 */
import (
	"github.com/hq-cml/spider-face/core"
	"github.com/hq-cml/spider-face"
	"fmt"
	"time"
)

func main() {
	//server config
	sConfig := &core.SpiderConfig{
		BindAddr: ":9529",    //监听地址:端口
	}

	//生成Spider实例
	spd, err := spider.NewSpider(sConfig, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	//快捷注册样例
	spd.GET("/", func(rp core.Roundtrip) {
		rp.Echo("Hello world!")
	})

	spd.GET("/api", func(rp core.Roundtrip) {             //普通的API接口输出
		rp.Echo("Get参数Id是：" + rp.Param("id"))
	})

	spd.GET("/quick/display", func(rp core.Roundtrip) {   //展示一个页面（基于模板文件）
		rp.Assign("nowtime", time.Now())                     //给模板传参
		rp.Assign("title", "welcome to spider~")
		rp.Assign("id", rp.Param("id"))

		rp.Display("index")                              //渲染展示index.html模板
	})

	spd.POST("/quick/post", func(rp core.Roundtrip) {     //Post接口
		rp.Echo("Post参数Id是：" + rp.Param("id"))
	})

	//Run
	spd.Run()
}
