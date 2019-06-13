# 开发文档

## 项目背景
Golang是21实际的C语言，它的显著风格就是清爽、直率，不弯弯绕，并且goang对http有着天然支持。使用golang，你只需要寥寥几行代码，就可以写出一个可以独立运行部署的web server的接口。不需要像PHP、Java那样搭配一套apache，tomcat服务器软件才能使用。  

没有恼人的依赖、不需要复杂的配置，一切都是顺其自然，并且，他还天然支持高并发（基于goroutine）。
这样的特点，都非常吸引人。我一直觉得互联网项目应该遵循：简单、皮实、可靠的原则。

## 框架的意义
Golang相关的web框架已经有很多了，像echo、gin、beego、martini，自己实现一个框架的意义何在？  
我认为自己搞一个框架可以更好的理解golang对于http的设计与理解。并且对于一些特殊功能，也更加可控。
  
如前所述，golang对web开发已经支持的很到位，但是，我个人认为，即便如此，还是有一定门槛。
- 首先，用户必须对http,tcp等网络协议多多少少有一些了解。
- 其次，对于go实现http web的套路，也需要有一个熟悉和理解的过程。这些都是门槛。

基于这样事实，我个人认为框架的意义在于：**进一步降低门槛。我希望Spider-face有以下的功效：**  
- 进一步屏蔽http协议细节，封装统一的接口，提供便捷的操作函数实现常用的http控制。
- 进一步屏蔽go实现http的细节，用户只需要会go语言本身语法，就可以实现一个http service。
- 提供好用的路由管理包括rewrite功能（这一部分是代替传统的web server如apache、nginx等）。
- 其他细节的屏蔽，总之就是希望用户将精力集中在业务逻辑本身，无感不相干的底层技术栈干扰。



## 通用概念解释
MVC：
Action：
路由：
//TODO

## 使用方法

### 配置文件解析
配置结构体，可以用于控制spider的行为，作为参数传入spider的创建函数中
```
type SpiderConfig struct {
	BindAddr      string //程序绑定的地址：默认地址:9529
	TplPath       string //模板文件存放位置：默认位置为可执行程序目录下的/tpl
	StaticPath    string //静态文件根目录：默认位置为可执行程序目录下的/static
	LogPath       string //日志位置：默认直接输出到终端上
	LogLevel      string //日志级别：Debug/Info/Error/Fatal。默认是Debug
	
	ReadTimeout   		int64  //服务器读取request的总时间
	WriteTimeout  		int64  //服务器回写总时间
	MaxHeaderByte 		int64  //最大的Header的大小

	CustomHttpErrorHtml map[int]string    //定制化的错误页面 httpCode => customErr.html
	CustomRewriteRule   map[string]string //rewrite规则，见下面详述
}
```