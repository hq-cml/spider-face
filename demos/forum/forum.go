package main

import (
	"github.com/hq-cml/spider-face"
	"github.com/hq-cml/spider-face/core"
)

func main() {
	//配置
	sConfig := &core.SpiderConfig {
		BindAddr: ":9529",
	}

	//生成实例
	spd := spider.NewSpider(sConfig, nil)

	//Run
	spd.Run()
}