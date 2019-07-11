package main

/*
* 利用快捷注册方式
* 快速开始一个web项目
* 如果是一个简单的项目，提供简单的页面，快捷注册是合适的选择
*
* 如果项目比较复杂，需要合理的组织结构，请使用“通用模式”
*/
import (
	"github.com/hq-cml/spider-face/core"
	"github.com/hq-cml/spider-face"
	"time"
)

func main() {
	spd := spider.NewSpider(nil, nil)           //生成Spider实例, 默认地址:端口

	spd.GET("/index", func(rp core.Roundtrip) {       //展示一个页面（基于模板文件）
		rp.Assign("nowtime", time.Now())                  //给模板传参
		rp.Assign("title", "welcome to spider~")
		rp.Assign("id", rp.Param("id"))              //接收一个参数，然后传给模板

		rp.Display("index")                           //渲染展示index.html模板
	})

	//Run
	spd.Run()
}
