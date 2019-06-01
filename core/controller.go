package core

import (
	"mime/multipart"
	"net/http"
	"fmt"
)

/*
 * Controller接口，规定了Spider中必需的合法的行为
 */
type Controller interface {
	//用户实现
	GetRouter() []ControllerRouter

	//Spider提供, 获取Controller的实时状态用于任务分发场景
	Init(request *Request, response *Response, logger SpiderLogger) bool
	SetName(name string)
	GetName() string
	SetAction(name string)
	GetAction() string
	Echo(content string)
	Param(key string, defaultValue ...string) string
	Display(viewPath ...string)
	Assign(key interface{}, value interface{})

}

//实时controller，每次请求到来，都会动态生成一个controller实例
//通过这个实例来控制逻辑处理、输入、输出
//这个结构将会嵌入到所有的用户定制Controller对象中去
type RuntimeController struct {
	request  *Request
	response *Response
	view     *View

	controllerName string
	actionName     string
}

type ControllerRouter struct {
	Method   string
	Location string
	Action   string
}

func (rc *RuntimeController) Init(request *Request, response *Response, logger SpiderLogger) bool {
	rc.request = request
	rc.response = response
	rc.view = NewView(logger)

	return true
}

func (rc *RuntimeController) SetName(name string) {
	if name == "" {
		return
	}
	rc.controllerName = name
}

func (rc *RuntimeController) GetName() string {
	return rc.controllerName
}

func (rc *RuntimeController) SetAction(name string) {
	if name == "" {
		return
	}
	rc.actionName = name
}

func (rc *RuntimeController) GetAction() string {
	return rc.actionName
}

func (rc *RuntimeController) Param(key string, defaultValue ...string) string {
	v := rc.request.FindParam(key)
	if v == "" && defaultValue != nil {
		return defaultValue[0]
	}
	return v
}

//向页面模板引擎注册数据,待展示用
func (rc *RuntimeController) Assign(key interface{}, value interface{}) {
	rc.view.Assign(key, value)
}

//输出展示页面
func (rc *RuntimeController) Display(viewPath ...string) {
	bytes, err := rc.Render(viewPath...)

	if err == nil {
		rc.response.SetHeader("Content-Type", "text/html; charset=utf-8")
		rc.response.WriteBody(bytes)
	} else {
		rc.response.SetHeader("Content-Type", "text/html; charset=utf-8")
		rc.response.WriteBody([]byte(err.Error()))
	}
}

func (rc *RuntimeController) Render(viewPath ...string) ([]byte, error) {
	var viewPathName string
	if viewPath == nil || viewPath[0] == "" {
		viewPathName = rc.GetName() + "/" + rc.GetAction()
		fmt.Println("viewName:", viewPathName)
	} else {
		viewPathName = viewPath[0]
		fmt.Println("viewName_x:", viewPathName)
	}
	return rc.view.Render(viewPathName)
}

func (rc *RuntimeController) GetCookie(name string) string {
	return rc.request.GetCookie(name)
}

func (rc *RuntimeController) GetUri() string {
	return rc.request.GetUri()
}

func (rc *RuntimeController) UrlPath() string {
	return rc.request.UrlPath()
}

func (rc *RuntimeController) GetClientIP() string {
	return rc.request.GetClientIP()
}

func (rc *RuntimeController) Scheme() string {
	return rc.request.Scheme()
}

func (rc *RuntimeController) Header(key string) string {
	return rc.request.GetHeader(key)
}

func (rc *RuntimeController) SetHeader(key, value string) {
	rc.response.SetHeader(key, value)
}

func (rc *RuntimeController) SetCookie(name string, value string, others ...interface{}) {
	rc.response.SetCookie(name, value, others...)
}

func (rc *RuntimeController) Echo(content string) {
	rc.OutputBytes([]byte(content))
}

func (rc *RuntimeController) OutputBytes(bytes []byte) {
	rc.response.SetHeader("Content-Type", "text/html; charset=utf-8")
	rc.response.WriteBody(bytes)
}

func (rc *RuntimeController) OutputJson(data interface{}, coding ...bool) error {
	return rc.response.Json(data, coding...)
}

func (rc *RuntimeController) OutputJsonp(callback string, data interface{}, coding ...bool) error {
	return rc.response.Jsonp(callback, data, coding...)
}

func (rc *RuntimeController) GetMethod() string {
	return rc.request.GetMethod()
}

//获取所有get变量
func (rc *RuntimeController) GET() map[string]string {
	return rc.request.GetAllGetParams()
}

//获取所有post提交变量
func (rc *RuntimeController) POST() map[string]interface{} {
	return rc.request.GetAllPostParams()
}

//获取request的body
func (rc *RuntimeController) ReqBody() []byte {
	return rc.request.ReadBody()
}

//跳转
func (rc *RuntimeController) Redirect(url string) {
	http.Redirect(rc.response.Writer, rc.request.request, url, 301)
}

//TODO
//获取上传文件
func (rc *RuntimeController) GetUploadFiles(key string) ([]*multipart.FileHeader, error) {
	return rc.request.GetUploadFiles(key)
}

func (rc *RuntimeController) MoveUploadFile(fromfile, tofile string) error {
	return rc.request.MoveUploadFile(fromfile, tofile)
}

func (rc *RuntimeController) GetFileSize(file *multipart.File) int64 {
	return rc.request.GetFileSize(file)
}
