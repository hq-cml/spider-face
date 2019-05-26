package core

import (
	"mime/multipart"
	"net/http"
	"fmt"
)

type BaseController struct {
	request  *Request
	response *Response
	view     *View
}

type SpiderController interface {
	GetRouter() map[string]interface{}
}

func (this *BaseController) Init(request *Request, response *Response) bool {
	this.request = request
	this.response = response
	this.view = NewView()

	return true
}

func (this *BaseController) Param(key string, defaultValue ...string) string {
	v := this.request.Param(key)
	if v == "" && defaultValue != nil {
		return defaultValue[0]
	}
	return v
}

func (this *BaseController) Assign(key interface{}, value interface{}) {
	this.view.Assign(key, value)
}

func (this *BaseController) Display(viewPath ...string) {
	bytes, err := this.Render(viewPath...)

	if err == nil {
		this.response.Header("Content-Type", "text/html; charset=utf-8")
		this.response.Body(bytes)
	} else {
		this.response.Header("Content-Type", "text/html; charset=utf-8")
		this.response.Body([]byte(err.Error()))
	}
}

func (this *BaseController) Render(viewPath ...string) ([]byte, error) {
	var view_name string
	if viewPath == nil || viewPath[0] == "" {
		view_name = this.request.GetController() + "/" + this.request.GetAction()
		fmt.Println("viewName:", view_name)
	} else {
		view_name = viewPath[0]
	}
	return this.view.Render(view_name)
}

func (this *BaseController) Cookie(name string) string {
	return this.request.Cookie(name)
}

func (this *BaseController) Uri() string {
	return this.request.Uri()
}

func (this *BaseController) UrlPath() string {
	return this.request.UrlPath()
}

func (this *BaseController) IP() string {
	return this.request.IP()
}

func (this *BaseController) Scheme() string {
	return this.request.Scheme()
}

func (this *BaseController) Header(key string) string {
	return this.request.Header(key)
}

func (this *BaseController) SetHeader(key, value string) {
	this.response.Header(key, value)
}

func (this *BaseController) SetCookie(name string, value string, others ...interface{}) {
	this.response.Cookie(name, value, others...)
}

func (this *BaseController) Echo(content string) {
	this.OutputBytes([]byte(content))
}

func (this *BaseController) OutputBytes(bytes []byte) {
	this.response.Header("Content-Type", "text/html; charset=utf-8")
	this.response.Body(bytes)
}

func (this *BaseController) Json(data interface{}, coding ...bool) error {
	return this.response.Json(data, coding...)
}

func (this *BaseController) Jsonp(callback string, data interface{}, coding ...bool) error {
	return this.response.Jsonp(callback, data, coding...)
}

func (this *BaseController) Method() string {
	return this.request.Method()
}

//获取所有get变量
func (this *BaseController) GET() map[string]string {
	return this.request.ParamGet()
}

//获取所有post提交变量
func (this *BaseController) POST() map[string]interface{} {
	return this.request.ParamPost()
}

func (this *BaseController) Location(url string) {
	http.Redirect(this.response.Writer, this.request.request, url, 301)
}

//TODO
// 获取所有上传文件
// files, _ := this.GetUploadFiles("user_icon")
// for i, _ := range files {
//	 file, _ := files[i].Open()
//	 defer file.Close()
//	 log.Print(this.GetFileSize(&file))
// }
func (this *BaseController) GetUploadFiles(key string) ([]*multipart.FileHeader, error) {
	return this.request.GetUploadFiles(key)
}

func (this *BaseController) MoveUploadFile(fromfile, tofile string) error {
	return this.request.MoveUploadFile(fromfile, tofile)
}

func (this *BaseController) GetFileSize(file *multipart.File) int64 {
	return this.request.GetFileSize(file)
}
