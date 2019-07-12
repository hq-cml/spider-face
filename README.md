![标题](./img/SPIDER-FACE.png)

# 一个简单的web框架
Spider-Face是一个基于golang实现的web开发框架。  
初衷是我想给我的迷你搜索引擎 [Spider-Engine](https://github.com/hq-cml/spider-engine) 开发一套Web界面，后来索性就把它也抽象独立成了一套web框架了。

Spider-Face的目标是让golang的web开发更加简单！详细的使用方式，见：[开发文档](./design.md)

### 项目特点
* 基于经典的MVC模式（M这一块暂时还没实现，它与V,C相对独立，用户可以任意找一个熟悉的orm进行替换，比如 [xorm](http://www.xorm.io)）
* 支持快捷路由注册（这一点参考了echo框架，适用于快速开发一些只需要简单接口的项目）
* 支持传统的Controller/Action的路由注册（这种模式将按照MVC的模式来组织项目，适用于需要合理刮花结构的大型项目）
* 支持rewrite，类似于Nginx的rewrite配置的功能
* 支持更加简单的视图模板规则（使用户不需要理解golang的http.template包）

### 安装
```
$go get github.com/hq-cml/spider-face
```

### 第一个例子： Hello, World!

```
package main

import (
	"github.com/hq-cml/spider-face"
	"github.com/hq-cml/spider-face/core"
)

func main() {
	spd := spider.NewSpider(nil, nil)            //spd实例生成

	spd.GET("/index", func(rp core.Roundtrip) {  //快捷注册路由函数，一个"hello world"接口，诞生
		rp.Echo("Hello World!")
	})

	spd.Run()                                    //Run
}
```

### 启动服务
```
$ go run main.go
```

用浏览器访问 http://localhost:9529/index  
就能看到熟悉的， Hello World!

### Demos
[demos](./demos)