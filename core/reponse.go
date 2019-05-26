package core


import (
	"bytes"
	//"compress/flate"
	//"compress/gzip"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type Response struct {
	Writer  http.ResponseWriter
}

func NewResponse(w http.ResponseWriter) *Response {
	response := &Response{
		Writer:        w,
	}
	return response
}

// Set http response header
func (resp *Response) SetHeader(key, val string) {
	resp.Writer.Header().Set(key, val)
}

// Set http response code
func (resp *Response) SetHttpCode(code int) {
	resp.Writer.WriteHeader(code)
}

func (resp *Response) Body(html_content []byte) {
	//TODO 支持压缩
	//accept_encoding := resp.request.Header("Accept-Encoding")
	//if CompressType != COMPRESS_CLOSE && len(html_content) >= CompressMinSize && accept_encoding != "" && (strings.Index(accept_encoding, "gzip") >= 0 || strings.Index(accept_encoding, "flate") >= 0) {
	//	switch CompressType {
	//	case COMPRESS_GZIP:
	//		resp.Header("Content-Encoding", "gzip")
	//		output_writer, _ := gzip.NewWriterLevel(resp.Writer, gzip.BestSpeed)
	//		defer output_writer.Close()
	//		output_writer.Write(html_content)
	//	case COMPRESS_FLATE:
	//		resp.Header("Content-Encoding", "deflate")
	//		output_writer, _ := flate.NewWriter(resp.Writer, flate.BestSpeed)
	//		defer output_writer.Close()
	//		output_writer.Write(html_content)
	//	}
	//} else {
		resp.Writer.Write(html_content)
	//}
}

// Set cookie
// Copy from beego @https://github.com/astaxie/beego
func (resp *Response) SetCookie(name string, value string, others ...interface{}) {
	cookieNameFilter := strings.NewReplacer("\n", "-", "\r", "-")
	cookieValueFilter := strings.NewReplacer("\n", " ", "\r", " ", ";", " ")

	var b bytes.Buffer
	fmt.Fprintf(&b, "%s=%s", cookieNameFilter.Replace(name), cookieValueFilter.Replace(value))
	if len(others) > 0 {
		switch v := others[0].(type) {
		case int, int64, int32:
			vv, _ := strconv.ParseInt(fmt.Sprintf("%v", v), 10, 64)
			if vv > 0 {
				fmt.Fprintf(&b, "; Max-Age=%d", v)
			} else if vv < 0 {
				fmt.Fprintf(&b, "; Max-Age=0")
			}
		}
	}
	// the settings below
	// Path, Domain, Secure, HttpOnly
	// can use nil skip set

	// default "/"
	if len(others) > 1 {
		if v, ok := others[1].(string); ok && len(v) > 0 {
			fmt.Fprintf(&b, "; Path=%s", cookieValueFilter.Replace(v))
		}
	} else {
		fmt.Fprintf(&b, "; Path=%s", "/")
	}

	if len(others) > 2 {
		if v, ok := others[2].(string); ok && len(v) > 0 {
			fmt.Fprintf(&b, "; Domain=%s", cookieValueFilter.Replace(v))
		}
	}

	if len(others) > 3 {
		var secure bool
		switch v := others[3].(type) {
		case bool:
			secure = v
		default:
			if others[3] != nil {
				secure = true
			}
		}
		if secure {
			fmt.Fprintf(&b, "; Secure")
		}
	}

	// default false. for session cookie default true
	httponly := false
	if len(others) > 4 {
		if v, ok := others[4].(bool); ok && v {
			// HttpOnly = true
			httponly = true
		}
	}

	if httponly {
		fmt.Fprintf(&b, "; HttpOnly")
	}

	resp.Writer.Header().Add("Set-Cookie", b.String())
}

// Set output type:json
func (resp *Response) Json(data interface{}, coding ...bool) error {
	resp.SetHeader("Content-Type", "application/json;charset=UTF-8")

	var content []byte
	var err error

	content, err = json.Marshal(data)
	if err != nil {
		http.Error(resp.Writer, err.Error(), http.StatusInternalServerError)
		return err
	}
	if coding != nil && coding[0] == true {
		content = []byte(unicode(string(content)))
	}
	resp.Body(content)
	return nil
}

func (resp *Response) Jsonp(callback string, data interface{}, coding ...bool) error {

	resp.SetHeader("Content-Type", "application/javascript;charset=UTF-8")

	var content []byte
	var err error

	content, err = json.Marshal(data)
	if err != nil {
		http.Error(resp.Writer, err.Error(), http.StatusInternalServerError)
		return err
	}
	if coding != nil && coding[0] == true {
		content = []byte(unicode(string(content)))
	}
	ck := bytes.NewBufferString(" " + callback)
	ck.WriteString("(")
	ck.Write(content)
	ck.WriteString(");\r\n")

	resp.Body(ck.Bytes())
	return nil
}

// Convert to unicode
// TODO 测试
func unicode(str string) string {
	rs := []rune(str)
	jsons := ""
	for _, r := range rs {
		rint := int(r)
		if rint < 128 {
			jsons += string(r)
		} else {
			jsons += "\\u" + strconv.FormatInt(int64(rint), 16) // json
		}
	}
	return jsons
}
