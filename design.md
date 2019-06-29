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
#### MVC
MVC是一个设计模式，它强制性的使应用程序的输入、处理和输出分开。使用MVC应用程序被分成三个核心部件：模型（M）、视图（V）、控制器（C），它们各自处理自己的任务。  
- 视图（V）：视图是用户看到并与之交互的界面。简单理解视图就是由HTML元素组成的界面，它只是作为一种输出数据并允许用户操纵的方式。 
- 模型（M）：实际负责任务的处理，包括数据的读取、组装等，模型与数据格式无关，这样一个模型能为多个视图提供数据。（PS，由于时间关系暂时还未实现模型，因为它与C和V相对独立，用户可以自行实现或者封装）
- 控制器（C）：控制器接受用户的输入，然后调用模型用业务逻辑来处理用户的请求并返回数据，通过视图去展示用户的需求。
- Action：即所谓的行为动作，它通常是控制器的一个方法，对外表现形式通常是为一个api接口
- 路由：将一个URI路径和一个控制器的某个Action进行绑定。当访问该URI的时候，框架将调用对应的Action方法进行处理。


PS：这里我只列了最重要的几个概念，其实周边概念还有不少，可以查看[ThinkPHP](http://document.thinkphp.cn/manual_3_2.html)（一个php的web框架）的文档。我基本上是以ThinkPHP的行为作为蓝本，设计并实现的Spider-Face。


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
	
	Mime                bool              //配置Mime
	Gzip                bool              //开启gzip压缩，开启后，默认对：.css、.js、.html、.jpg、.png进行压缩
	CustomGzipExt       string            //用户自己指定的压缩文件后缀，用|分隔，比如: jpg|css|js|png
}
```


### 路由种类：
Spider的特色是，支持两种路由注册方式，快捷注册和通用注册。
- 快捷注册：以最简化的方式快速实现一个接口逻辑并绑定一个URL路由。
- 通用注册：按照controller/action的组织结构在注册路由。


模式| 优点 | 不足 | 适用场景  |  不支持的功能
--- | ---|---|---|---
快捷模式 | 简单快速 | 部分高级功能不支持，组织结构不清晰 | 快速开发一些简单接口 | 路径参数，pathinfo参数，不常用的http method等
通用模式 | 支持高级功能，组织管理清晰   成本稍高 | 成本稍高 | 一个复杂的项目，需要多层管理| -


### 路由注册方法

#### 快捷模式
 仅支持最常用的四中Method，注册方式如下：
```
spd.POST("/users", core.ActionFunc)
spd.GET("/users", core.ActionFunc)
spd.PUT("/users", core.ActionFunc)
spd.DELETE("/users", core.ActionFunc)
```

其中，core.ActionFunc是函数类型：func(rp Roundtrip)，详情见：[例子](./demos/quick-echo)

#### 通用模式
通用模式适合需要复杂组织结构的大型项目，所以也就比快捷模式多一些约束。  
Spider遵从约定大于配置的原则，一些需要默认的约定如下：
- contrller的命名必须以字符**Controller**作为后缀，比如IndexController
- action的命名必须要字符**Action**作为后缀，比如DeleteAction
- controller的类定义，不强制目录位置，但是建议放在统一的controllers目录中，便于管理。
- 建议每个controller类，使用一个独立的文件存放。

说起来抽象，其实也不复杂，直接看： [例子](./demos/common-echo)

### 关于Roundtrip
Roundtrip是spider提供给用户的操作抓手，它表示一次完整请求与响应的。利用Roundtrip，用户可以实现：
- 参数获取
- 结果返回（适用于api）
- 页面返回（适用于web页面）

```
//具体Roundtrip提供的功能如下：
type Roundtrip interface {
	GetControllerName() string
	GetActionName() string
	GetRequest() *Request
	GetResponse() *Response
	Param(key string, defaultValue ...string) string
	Display(viewPath ...string)
	Assign(key interface{}, value interface{})
	Render(viewPath ...string) ([]byte, error)
	GetCookie(name string) string
	GetUri() string
	UrlPath() string
	GetClientIP() string
	Scheme() string
	Header(key string) string
	SetHeader(key, value string)
	SetCookie(name string, value string, others ...interface{})
	Echo(content string)
	OutputBytes(bytes []byte)
	OutputJson(data interface{}, coding ...bool) error
	OutputJsonp(callback string, data interface{}, coding ...bool) error
	GetMethod() string
	GET() map[string]string
	POST() map[string]interface{}
	ReqBody() []byte
	Redirect(url string, code ...int)
	GetUploadFiles(key string) ([]*multipart.FileHeader, error)
	MoveUploadFile(fromfile, tofile string) error
	GetFileSize(file *multipart.File) int64
}
```

#### 参数读取
rp.Param方法，兼容了POST参数和GET参数
此时假如访问:

```
curl -X GET 'http://127.0.0.1/index?id=123' 或者
curl -X POST 'http://127.0.0.1/index' -d 'id=123'，
```

通过rp.Param("id")，可以拿到id参数

#### POST的Body读取
有的时候Post参数不是类似application/x-www-form-urlencoded这样的form表单，比如。

```
curl -X POST 'http://127.0.0.1/index' -d '{"id":"123", "name":"haha"}'，
```

那么可以用如下形式将body读取，然后自行解码。如下：

```
var body []byte
body = rp.ReqBody()
```

#### 读取请求的Header


#### 返回结果
对于一些api接口，往往需要直接返回内容，可以如下

```
rp.Echo("Hello")
```

对于一些需要结构化的场景，往往需要返回json，可以直接使用如下快捷函数：

```
rp.Json(x) //x是一个结构的变量或指针
```

### 关于视图（模板）
golang本身对于template就有一定的支持（但是有一定的学习门槛）。Spider没有重新造轮子，在http.template的基础上，进行了一定的封装，力图降低这个门槛。  
所以相对应的，spider对于模板的使用，也有一些自己的规则（注意这不是golang的http.template规则，而是spider为了降低门槛封装后的规则）  
- 模板文件必须是.html后缀的文本文件
- 一个模板文件，只定义一个模板。
- 模板文件支持include操作，操作符为{{ template xxx }}
- 用{{ define xxx/xxx }}来定义模板，其中“xxx/xxx”必须严格按照视图模板文件目录的相对路径，且无需后缀名。
- 模板文件的名字，必须和“xxx/xxx”中的最后一段完全一致，不包括.html后缀

eg:  
假设目录/tmp/spider-ui/tpl/模板文件目录。那么模板文件/tmp/spider-ui/tpl/aa/bb.html，它的内容只能是

```
{{ define aa/bb }}
	...
{{ end }}
```

#### 输出页面
Spider采用了MVC的形式，所以如果需要输出页面，那么需要先定义模板  
定义好模板之后，可以如下方式使用模板：

```
rp.Assign("param1", "Hello")  //给模板赋值
rp.Assign("param2", "World")
rp.Display(templateName)      //输出模板，参数是template名字
```


如果你完全按照标准约定存放模板文件，那么你甚至可以直接使用

```
pr.Display() //程序会自动去./tpl/$controller/$action.html中查找文件。
```
$controller 和 $action对应于实际的控制器和动作名称（转lower且不带后缀）


##### 视图例子
说起来有点抽象，其实已经尽力将门槛降到go原生的模板之下了~~.  
[show you the code!](./demos/quick-html)


# 高级功能
### URL路径参数
路径参数仅适用于通用模式，它又分为两种形式：

```
普通形式：
	/user/:id => 访问/user/101等价于访问/user?id=101
	
PathInfo形式（直接以***结尾）：
	/orders/*** => 访问/orders/year/2019/month/06/day/05等价于访问/orders?year=2019&month=06&day=05
```

### rewrite功能
这个功能类似于nginx的rewrite重写。要使用rewrite功能，只需要在配置结构里面具体，比如：

```
sConfig.CustomRewriteRule = map[string]string {
    "/test/rewrite": "/index?name=123",
    "/test/rewrite/(.*)/(.*)": "/index?id=$1&name=$2",
}

//这样所有发向/test/rewrite地址的请求，将被重定向到/index上，并且还会带一个name=123参数。
//这样所有发向/test/rewrite/12/abc地址的请求，将被重定向到/index上，并且还会带一个id=12&name=abc参数。
```

可以看出，这是一种路径参数变体，也就是说rewrite功能也可以变相实现路径参数


### TODO
* 参数绑定，即对于传入的结构化的post参数，直接映射到一个类实例
* 插件功能
* 支持热升级:不中断服务重启server(抄袭beego ^_^)
* 支持自定义路由
* 支持https
* 支持gzip/deflate压缩
* 支持静态文件