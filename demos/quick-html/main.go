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
	"fmt"
)

func main() {
	spd := spider.NewSpider(nil, nil)           //生成Spider实例, 默认地址:端口

	spd.GET("/index", func(rp core.Roundtrip) {       //展示一个页面（基于模板文件）
		rp.Assign("nowtime", time.Now())                  //给模板传参
		rp.Assign("title", "welcome to spider~")
		rp.Assign("id", rp.Param("id"))              //接收一个参数，然后传给模板

		rp.Display("index")                           //渲染展示index.html模板
	})

	spd.POST("/upload", func(rp core.Roundtrip) {       //展示一个页面（基于模板文件）
		key := "uploadfile"
		mHeaders, err := rp.GetUploadFiles(key)
		if err != nil {
			rp.Echo("Some thing wrong!")
		}

		rp.Echo(fmt.Sprintf("成功上传文件个数：%d", len(mHeaders)))
		rp.Echo("<Br>")
		header := mHeaders[0]
		rp.Echo("文件名称：")
		rp.Echo(header.Filename)
		rp.Echo("<Br>")
		//rp.Echo("文件大小：")
		//rp.Echo(fmt.Sprintf("%d", rp.GetFileSize(header)))

		rp.MoveUploadFile(key, "/data/share/"+header.Filename)
		rp.Echo("上传完成")
	})

	//Run
	spd.Run()
}
