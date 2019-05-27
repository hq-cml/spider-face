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

	//Spider提供
	Init(request *Request, response *Response) bool
	SetController(name string)
	SetAction(name string)
	GetController() string
	GetAction() string
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
	Method string
	Pattern string
	Action string
}

func (rc *RuntimeController) Init(request *Request, response *Response) bool {
	rc.request = request
	rc.response = response
	rc.view = NewView()

	return true
}

func (rc *RuntimeController) SetController(name string) {
	if name == "" {
		return
	}
	rc.controllerName = name
}

func (rc *RuntimeController) SetAction(name string) {
	if name == "" {
		return
	}
	rc.actionName = name
}

func (rc *RuntimeController) GetController() string {
	return rc.controllerName
}

func (rc *RuntimeController) GetAction() string {
	return rc.actionName
}

func (rc *RuntimeController) Param(key string, defaultValue ...string) string {
	v := rc.request.Param(key)
	if v == "" && defaultValue != nil {
		return defaultValue[0]
	}
	return v
}

func (rc *RuntimeController) Assign(key interface{}, value interface{}) {
	rc.view.Assign(key, value)
}

func (rc *RuntimeController) Display(viewPath ...string) {
	bytes, err := rc.Render(viewPath...)

	if err == nil {
		rc.response.SetHeader("Content-Type", "text/html; charset=utf-8")
		rc.response.Body(bytes)
	} else {
		rc.response.SetHeader("Content-Type", "text/html; charset=utf-8")
		rc.response.Body([]byte(err.Error()))
	}
}

func (rc *RuntimeController) Render(viewPath ...string) ([]byte, error) {
	var view_name string
	if viewPath == nil || viewPath[0] == "" {
		view_name = rc.GetController() + "/" + rc.GetAction()
		fmt.Println("viewName:", view_name)
	} else {
		view_name = viewPath[0]
	}
	return rc.view.Render(view_name)
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
	rc.response.Body(bytes)
}

func (rc *RuntimeController) Json(data interface{}, coding ...bool) error {
	return rc.response.Json(data, coding...)
}

func (rc *RuntimeController) Jsonp(callback string, data interface{}, coding ...bool) error {
	return rc.response.Jsonp(callback, data, coding...)
}

func (rc *RuntimeController) GetMethod() string {
	return rc.request.GetMethod()
}

//获取所有get变量
func (rc *RuntimeController) GET() map[string]string {
	return rc.request.ParamGet()
}

//获取所有post提交变量
func (rc *RuntimeController) POST() map[string]interface{} {
	return rc.request.ParamPost()
}

//跳转
func (rc *RuntimeController) Location(url string) {
	http.Redirect(rc.response.Writer, rc.request.request, url, 301)
}

//TODO
// 获取所有上传文件
// files, _ := this.GetUploadFiles("user_icon")
// for i, _ := range files {
//	 file, _ := files[i].Open()
//	 defer file.Close()
//	 log.Print(this.GetFileSize(&file))
// }
func (rc *RuntimeController) GetUploadFiles(key string) ([]*multipart.FileHeader, error) {
	return rc.request.GetUploadFiles(key)
}

func (rc *RuntimeController) MoveUploadFile(fromfile, tofile string) error {
	return rc.request.MoveUploadFile(fromfile, tofile)
}

func (rc *RuntimeController) GetFileSize(file *multipart.File) int64 {
	return rc.request.GetFileSize(file)
}
