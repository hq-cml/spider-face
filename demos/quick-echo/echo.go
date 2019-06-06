package main

/*
 * 利用快捷注册方式
 * 注册一些API接口
 * 如果是一个简单的项目，提供http的接口，快捷注册是合适的选择
 *
 * 如果项目比较复杂，需要合理的组织结构，请使用“通用模式”
 */
import (
	"github.com/hq-cml/spider-face"
	"github.com/hq-cml/spider-face/core"
)

func main() {
	spd, err := spider.NewSpider(&core.SpiderConfig{    //生成Spider实例
		BindAddr: ":9529",    						    //监听地址:端口
	}, nil)
	if err != nil {
		return
	}

	spd.GET("/index", func(rp core.Roundtrip) {       //快捷注册路由函数，一个"hello world"接口，诞生
		rp.Echo("Hello World!")
	})

	spd.GET("/api/get", func(rp core.Roundtrip) {     //Get接口接收参数：curl 'http://ip:9529/api/get?id=aaa'
		rp.Echo("Get参数Id是：" + rp.Param("id"))
	})

	spd.POST("/api/post", func(rp core.Roundtrip) {   //Post接口接收参数：curl -X POST 'http://ip:9529/api/post' -d 'id=aaa'
		rp.Echo("Post参数Id是：" + rp.Param("id"))
	})

	spd.GET("/json", func(rp core.Roundtrip) {        //快捷注册路由函数，用json输出一个结构化的数据
		m := map[string]interface{} {
			"A": "a",
			"B": "b",
			"C": []int{1,2,3},
		}
		rp.OutputJson(m)
	})

	spd.Run()                                                  //Run
}
