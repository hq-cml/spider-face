package core

import (
	"os"
	"net/http"
	"fmt"
)

//禁用默认的IE 和chrome错误页面显示
var disableIEAndChrome = `
<!-- a padding to disable MSIE and Chrome friendly error page -->
<!-- a padding to disable MSIE and Chrome friendly error page -->
<!-- a padding to disable MSIE and Chrome friendly error page -->
<!-- a padding to disable MSIE and Chrome friendly error page -->
<!-- a padding to disable MSIE and Chrome friendly error page -->
<!-- a padding to disable MSIE and Chrome friendly error page -->`

var ErrorPagesMap = map[int]string {
	403 : `<html>
<head><title>403 Forbidden</title></head>
<body bgcolor="white">
<center><h1>403 Forbidden</h1></center>
<hr><center>spider/0.0.1</center>
</body>
</html>`,

	404 :`<html>
<head><title>404 Not Found</title></head>
<body bgcolor="white">
<center><h1>404 Not Found</h1></center>
<hr><center>spider/0.0.1</center>
</body>
</html>`,

	500 : `<html>
<head><title>500 Internal Server Error</title></head>
<body bgcolor="white">
<center><h1>500 Internal Server Error</h1></center>
<hr><center>spider/0.0.1</center>
</body>
</html>`,
}

func OutputStaticFile(response *Response, request *Request, file string, customErrHtml map[int]string) {
	filePath := GlobalConf.StaticPath + file
	fi, err := os.Stat(filePath)
	if err != nil && os.IsNotExist(err) {
		OutputErrorHtml(response, request, http.StatusNotFound, customErrHtml)
		return
	} else if fi.IsDir() == true {
		OutputErrorHtml(response, request, http.StatusForbidden, customErrHtml)
		return
	}
	//file_size := fi.Size()
	//mod_time := fi.ModTime()

	http.ServeFile(response.Writer, request.request, filePath)
	return
	//TODO 压缩

}

func OutputErrorHtml(response *Response, request *Request, httpCode int, customErrHtml map[int]string) {
	//用户自定义的错误页面
	if customErrHtml != nil {
		if errHtml, exist := customErrHtml[httpCode]; exist {
			if fi, err := os.Stat(errHtml); (err == nil || os.IsExist(err)) && fi.IsDir() != true {
				http.ServeFile(response.Writer, request.request, errHtml)
				return
			}
		}
	}

	//设置HTTP Repsonse Header
	response.SetHeader("Status", fmt.Sprintf("%d", httpCode))
	response.SetHeader("Content-Type", "text/html; charset=utf-8")
	response.SetHeader("X-Content-Type-Options", "nosniff")

	//设置HTTP CODE
	response.SetHttpCode(httpCode)

	//回写HTTP Body
	fmt.Fprintln(response.Writer, ErrorPagesMap[httpCode] + disableIEAndChrome)
}
