package core


import (
	"io"
	"os"
	"errors"
	"net/http"
	"strings"
	"mime/multipart"
)

type Request struct {
	request    *http.Request
	formParsed bool
	pathParams map[string]string  //保存通过URL路径传入的参数
}

func NewRequest(r *http.Request) *Request {
	request_instance := &Request{
		request: r,
	}

	return request_instance
}

func (req *Request) GetMethod() string {
	return req.request.Method
}

func (req *Request) GetUri() string {
	return req.request.RequestURI
}

func (req *Request) UrlPath() string {
	return req.request.URL.Path
}

func (req *Request) GetClientIP() string {
	ips := req.Proxy()
	if len(ips) > 0 && ips[0] != "" {
		rip := strings.Split(ips[0], ":")
		return rip[0]
	}
	ip := strings.Split(req.request.RemoteAddr, ":")
	if len(ip) > 0 {
		if ip[0] != "[" {
			return ip[0]
		}
	}
	return "127.0.0.1"
}

func (req *Request) Proxy() []string {
	if ips := req.GetHeader("X-Forwarded-For"); ips != "" {
		return strings.Split(ips, ",")
	}
	return []string{}
}

func (req *Request) Scheme() string {
	if req.request.URL.Scheme != "" {
		return req.request.URL.Scheme
	}
	if req.request.TLS == nil {
		return "http"
	}
	return "https"
}

func (req *Request) GetHeader(key string) string {
	return req.request.Header.Get(key)
}

// Get cookie
func (req *Request) GetCookie(key string) string {
	cookie, err := req.request.Cookie(key)
	if err != nil {
		return ""
	}
	return cookie.Value
}

//对参数就进行语法分析
func (req *Request) parseParam() error {
	if req.formParsed == true || req.request.Form != nil || req.request.PostForm != nil || req.request.MultipartForm != nil {
		return nil
	}

	req.formParsed = true
	if strings.Contains(req.GetHeader("Content-Type"), "multipart/form-data") {
		if err := req.request.ParseMultipartForm(32 << 20); err != nil {
			return errors.New("Error parsing request body:" + err.Error())
		}
	} else if err := req.request.ParseForm(); err != nil {
		return errors.New("Error parsing request body:" + err.Error())
	}
	return nil
}

// Get request param by key
func (req *Request) Param(key string) string {
	if req.pathParams != nil {
		if d, exist := req.pathParams[key]; exist {
			return d
		}
	}
	err := req.parseParam()
	if err != nil {
		return ""
	}
	return req.request.Form.Get(key)
}

//TODO 获取POST的指定参数

//TODO 获取Body

// Get all request params passed by GET Method
func (req *Request) ParamGet() (data map[string]string) {
	err := req.parseParam()
	if err != nil {
		return nil
	}

	if req.request.Form == nil {
		return req.pathParams
	}
	data = make(map[string]string)
	for k, v := range req.request.Form {
		data[k] = v[0] //只取第一个
	}

	//追加上路径参数
	if req.pathParams != nil {
		for k, v := range req.pathParams {
			data[k] = v
		}
	}
	return data
}

// Get all request params passed by POST Method
func (req *Request) ParamPost() (data map[string]interface{}) {
	err := req.parseParam()
	if err != nil {
		return nil
	}
	if req.request.PostForm == nil {
		return nil
	}
	data = make(map[string]interface{})
	for k, v := range req.request.PostForm {
		if len(v) > 1 {
			data[k] = v
		} else {
			data[k] = v[0]
		}
	}
	return data
}


// Get all upload files
func (req *Request) GetUploadFiles(key string) ([]*multipart.FileHeader, error) {
	req.parseParam()

	if req.request.MultipartForm == nil {
		return nil, nil
	}
	files, ok := req.request.MultipartForm.File[key]
	if ok {
		return files, nil
	}
	return nil, http.ErrMissingFile
}

// Save upload file
func (req *Request) MoveUploadFile(fromfile, tofile string) error {
	file, _, err := req.request.FormFile(fromfile)
	if err != nil {
		return err
	}

	defer file.Close()

	f, err := os.OpenFile(tofile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer f.Close()
	io.Copy(f, file)
	return nil
}

type Size interface {
	Size() int64
}

// Get upload file size
func (req *Request) GetFileSize(file *multipart.File) int64 {
	if sizeInterface, ok := (*file).(Size); ok {
		return sizeInterface.Size()
	}
	return -1
}
