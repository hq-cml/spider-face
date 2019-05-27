package core


import (
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
)

type Request struct {
	request        *http.Request
	formParsed     bool
	rewriteParams  map[string]string
}

func NewRequest(r *http.Request) *Request {
	request_instance := &Request{
		request: r,
	}

	return request_instance
}

func (this *Request) Method() string {
	return this.request.Method
}

func (this *Request) Uri() string {
	return this.request.RequestURI
}

func (this *Request) UrlPath() string {
	return this.request.URL.Path
}

func (this *Request) IP() string {
	ips := this.Proxy()
	if len(ips) > 0 && ips[0] != "" {
		rip := strings.Split(ips[0], ":")
		return rip[0]
	}
	ip := strings.Split(this.request.RemoteAddr, ":")
	if len(ip) > 0 {
		if ip[0] != "[" {
			return ip[0]
		}
	}
	return "127.0.0.1"
}

func (this *Request) Scheme() string {
	if this.request.URL.Scheme != "" {
		return this.request.URL.Scheme
	}
	if this.request.TLS == nil {
		return "http"
	}
	return "https"
}

func (this *Request) Proxy() []string {
	if ips := this.Header("X-Forwarded-For"); ips != "" {
		return strings.Split(ips, ",")
	}
	return []string{}
}

func (this *Request) Header(key string) string {
	return this.request.Header.Get(key)
}

// Get cookie
func (this *Request) Cookie(key string) string {
	cookie, err := this.request.Cookie(key)
	if err != nil {
		return ""
	}
	return cookie.Value
}

// Get request param by key
func (this *Request) Param(key string) string {
	if this.rewriteParams != nil {
		if d, ok := this.rewriteParams[key]; ok == true {
			return d
		}
	}
	err := this.ParseMultiForm()
	if err != nil {
		return ""
	}
	return this.request.Form.Get(key)
}

// Get all request params passed by GET Method
func (this *Request) ParamGet() (data map[string]string) {
	err := this.ParseMultiForm()
	if err != nil {
		return nil
	}

	if this.request.Form == nil {
		return this.rewriteParams
	}
	data = make(map[string]string)
	for k, v := range this.request.Form {
		data[k] = v[0]
	}
	if this.rewriteParams != nil {
		for k, v := range this.rewriteParams {
			data[k] = v
		}
	}
	return data
}

// Get all request params passed by POST Method
func (this *Request) ParamPost() (data map[string]interface{}) {
	err := this.ParseMultiForm()
	if err != nil {
		return nil
	}
	if this.request.PostForm == nil {
		return nil
	}
	data = make(map[string]interface{})
	for k, v := range this.request.PostForm {
		if len(v) > 1 {
			data[k] = v
		} else {
			data[k] = v[0]
		}
	}
	return data
}

func (this *Request) ParseMultiForm() error {
	if this.formParsed == true || this.request.Form != nil || this.request.PostForm != nil || this.request.MultipartForm != nil {
		return nil
	}

	this.formParsed = true
	if strings.Contains(this.Header("Content-Type"), "multipart/form-data") {
		if err := this.request.ParseMultipartForm(32 << 20); err != nil {
			return errors.New("Error parsing request body:" + err.Error())
		}
	} else if err := this.request.ParseForm(); err != nil {
		return errors.New("Error parsing request body:" + err.Error())
	}
	return nil
}

// Get all upload files
func (this *Request) GetUploadFiles(key string) ([]*multipart.FileHeader, error) {
	this.ParseMultiForm()

	if this.request.MultipartForm == nil {
		return nil, nil
	}
	files, ok := this.request.MultipartForm.File[key]
	if ok {
		return files, nil
	}
	return nil, http.ErrMissingFile
}

// Save upload file
func (this *Request) MoveUploadFile(fromfile, tofile string) error {
	file, _, err := this.request.FormFile(fromfile)
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
func (this *Request) GetFileSize(file *multipart.File) int64 {
	if sizeInterface, ok := (*file).(Size); ok {
		return sizeInterface.Size()
	}
	return -1
}
