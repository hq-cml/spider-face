package main

/*
 * 利用通用注册方式
 * 注册一个标准的controller，该controller拥有4个路由规则，每个规则都对应有Action（接口逻辑）
 * 并且，这4个规则拥不同的参数接收方式
 *
 * 如果是一个相对大型的项目，需要考虑组织结构，提供http的接口，通用注册是合适的选择
 */
import (
	"github.com/hq-cml/spider-face/core"
	"github.com/hq-cml/spider-face"
	"fmt"
	"time"
	"net/http"
)

//标准的Controller，必须以Controller后缀结尾
type HelloController struct {
	spdrp core.SpiderRoundtrip
}

//实现core.Controller接口
func (hello *HelloController) GetAllRouters() []core.ControllerRouter { //TODO 这种方式不够友好
	return []core.ControllerRouter{
		{Method: http.MethodGet, Location: "/index", Action:"IndexAction",},      //能接收普通参数
		{Method: http.MethodGet, Location: "/index/:id", Action:"IndexAction",},  //能接收普通参数和路径参数
		{Method: http.MethodGet, Location: "/index/*", Action:"IndexAction",},    //能接收普通参数和pathinfo参数
		{Method: http.MethodPost, Location: "/index/post", Action:"PostAction",}, //能接Post参数
	}
}
func (hello *HelloController) GetRoundTrip() core.Roundtrip {
	return &hello.spdrp
}

//curl 'http://192.168.110.133:9529/index/aaa'
//curl 'http://192.168.110.133:9529/index/name/hq/age/28'
func (hello *HelloController) IndexAction(rp core.Roundtrip) {
	rp.Echo(fmt.Sprintf("当前时间: %v\n", time.Now()))
	rp.Echo(fmt.Sprintf("参数Id: %v\n", rp.Param("id")))
	rp.Echo(fmt.Sprintf("参数name: %v\n", rp.Param("name")))
	rp.Echo(fmt.Sprintf("参数age: %v\n", rp.Param("age")))
}

//curl -X POST 'http://192.168.110.133:9529/index/post' -d 'id=aaa'
func (hello *HelloController) PostAction(rp core.Roundtrip) {
	rp.Echo(fmt.Sprintf("参数Id: %v\n", rp.Param("id")))
}

func main() {
	spd, err := spider.NewSpider(&core.SpiderConfig{    //生成Spider实例
		BindAddr: ":9529",    						    //监听地址:端口
	}, nil)
	if err != nil {
		return
	}

	if err = spd.RegisterController([]core.Controller{  //注册controller
		&HelloController{},
	}); err != nil {
		fmt.Println(err)
		return
	}

	spd.Run()                                           //Run
}
