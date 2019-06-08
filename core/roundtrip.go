package core

import (
	"mime/multipart"
	"net/http"
)

type Roundtrip interface {
	initRoundtrip(request *Request, response *Response, controllerName, actionName string, logger SpiderLogger) bool
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

//实时controller，每次请求到来，都会动态生成一个controller实例
//通过这个实例来控制逻辑处理、输入、输出
//这个结构将会嵌入到所有的用户定制Controller对象中去
type SpiderRoundtrip struct {
	request  *Request
	response *Response
	view     *View

	controllerName string
	actionName     string
}

func (rp *SpiderRoundtrip) initRoundtrip(request *Request, response *Response,
	controllerName, actionName string, logger SpiderLogger) bool {
	rp.request = request
	rp.response = response
	rp.controllerName = controllerName
	rp.actionName = actionName
	rp.view = NewView(logger)

	return true
}

func (rp *SpiderRoundtrip) GetResponse() *Response {
	return rp.response
}

func (rp *SpiderRoundtrip) GetRequest() *Request {
	return rp.request
}

func (rp *SpiderRoundtrip) GetControllerName() string {
	return rp.controllerName
}

func (rp *SpiderRoundtrip) GetActionName() string {
	return rp.actionName
}

func (rp *SpiderRoundtrip) Param(key string, defaultValue ...string) string {
	v := rp.request.FindParam(key)
	if v == "" && defaultValue != nil {
		return defaultValue[0]
	}
	return v
}

//向页面模板引擎注册数据,待展示用
func (rp *SpiderRoundtrip) Assign(key interface{}, value interface{}) {
	rp.view.Assign(key, value)
}

//输出展示页面
func (rp *SpiderRoundtrip) Display(viewPath ...string) {
	bytes, err := rp.Render(viewPath...)

	if err == nil {
		rp.response.SetHeader("Content-Type", "text/html; charset=utf-8")
		rp.response.WriteBody(bytes)
	} else {
		rp.response.SetHeader("Content-Type", "text/html; charset=utf-8")
		rp.response.WriteBody([]byte(err.Error()))
	}
}

func (rp *SpiderRoundtrip) Render(viewPath ...string) ([]byte, error) {
	var viewPathName string
	if viewPath == nil || viewPath[0] == "" {
		viewPathName = rp.GetControllerName() + "/" + rp.GetActionName()
	} else {
		viewPathName = viewPath[0]
	}
	return rp.view.Render(viewPathName)
}

func (rp *SpiderRoundtrip) GetCookie(name string) string {
	return rp.request.GetCookie(name)
}

func (rp *SpiderRoundtrip) GetUri() string {
	return rp.request.GetUri()
}

func (rp *SpiderRoundtrip) UrlPath() string {
	return rp.request.UrlPath()
}

func (rp *SpiderRoundtrip) GetClientIP() string {
	return rp.request.GetClientIP()
}

func (rp *SpiderRoundtrip) Scheme() string {
	return rp.request.Scheme()
}

func (rp *SpiderRoundtrip) Header(key string) string {
	return rp.request.GetHeader(key)
}

func (rp *SpiderRoundtrip) SetHeader(key, value string) {
	rp.response.SetHeader(key, value)
}

func (rp *SpiderRoundtrip) SetCookie(name string, value string, others ...interface{}) {
	rp.response.SetCookie(name, value, others...)
}

func (rp *SpiderRoundtrip) Echo(content string) {
	rp.OutputBytes([]byte(content))
}

func (rp *SpiderRoundtrip) OutputBytes(bytes []byte) {
	rp.response.SetHeader("Content-Type", "text/html; charset=utf-8")
	rp.response.WriteBody(bytes)
}

func (rp *SpiderRoundtrip) OutputJson(data interface{}, coding ...bool) error {
	return rp.response.Json(data, coding...)
}

func (rp *SpiderRoundtrip) OutputJsonp(callback string, data interface{}, coding ...bool) error {
	return rp.response.Jsonp(callback, data, coding...)
}

func (rp *SpiderRoundtrip) GetMethod() string {
	return rp.request.GetMethod()
}

//获取所有get变量
func (rp *SpiderRoundtrip) GET() map[string]string {
	return rp.request.GetAllGetParams()
}

//获取所有post提交变量
func (rp *SpiderRoundtrip) POST() map[string]interface{} {
	return rp.request.GetAllPostParams()
}

//获取request的body
func (rp *SpiderRoundtrip) ReqBody() []byte {
	return rp.request.ReadBody()
}

//跳转
func (rp *SpiderRoundtrip) Redirect(url string, code ...int) {
	if len(code) > 0 {
		http.Redirect(rp.response.Writer, rp.request.request, url, code[0])
	} else {
		http.Redirect(rp.response.Writer, rp.request.request, url, 302)
	}
}

//TODO
//获取上传文件
func (rp *SpiderRoundtrip) GetUploadFiles(key string) ([]*multipart.FileHeader, error) {
	return rp.request.GetUploadFiles(key)
}

func (rp *SpiderRoundtrip) MoveUploadFile(fromfile, tofile string) error {
	return rp.request.MoveUploadFile(fromfile, tofile)
}

func (rp *SpiderRoundtrip) GetFileSize(file *multipart.File) int64 {
	return rp.request.GetFileSize(file)
}
