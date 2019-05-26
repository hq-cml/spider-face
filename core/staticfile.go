package core

import (
	"os"
	"net/http"
	"fmt"
)

var ErrorPagesMap = map[int]string {
	403 :
	`<html>
	<head><title>403 Forbidden</title></head>
	<body bgcolor="white">
	<center><h1>403 Forbidden</h1></center>
	<hr><center>spider/0.0.1</center>
	</body>
	</html>
	`,
	404 :
	`
	<html>
	<head><title>404 Not Found</title></head>
	<body bgcolor="white">
	<center><h1>404 Not Found</h1></center>
	<hr><center>spider/0.0.1</center>
	</body>
	</html>
	`,
	500 :
	`
	<html>
	<head><title>500 Internal Server Error</title></head>
	<body bgcolor="white">
	<center><h1>500 Internal Server Error</h1></center>
	<hr><center>spider/0.0.1</center>
	</body>
	</html>
	<!-- a padding to disable MSIE and Chrome friendly error page -->
	<!-- a padding to disable MSIE and Chrome friendly error page -->
	<!-- a padding to disable MSIE and Chrome friendly error page -->
	<!-- a padding to disable MSIE and Chrome friendly error page -->
	<!-- a padding to disable MSIE and Chrome friendly error page -->
	<!-- a padding to disable MSIE and Chrome friendly error page -->
	`,
}

func OutputStaticFile(response *Response, request *Request, file string) {
	filePath := GlobalConf.StaticPath + file
	fi, err := os.Stat(filePath)
	if err != nil && os.IsNotExist(err) {
		OutErrorHtml(response, request, http.StatusNotFound)
		return
	} else if fi.IsDir() == true {
		OutErrorHtml(response, request, http.StatusForbidden)
		return
	}
	//file_size := fi.Size()
	//mod_time := fi.ModTime()

	http.ServeFile(response.Writer, request.request, filePath)

	//TODO 压缩

}

func OutErrorHtml(response *Response, request *Request, httpCode int) {
	//TODO
	//用户自定义的错误页面
	//if err_html, ok := response.server_config.HttpErrorHtml[httpCode]; ok == true {
	//	if fi, err := os.Stat(err_html); (err == nil || os.IsExist(err)) && fi.IsDir() != true {
	//		http.ServeFile(response.Writer, request.request, err_html)
	//		return
	//	}
	//}

	//设置HTTP Repsonse Header
	response.SetHeader("Status", fmt.Sprintf("%d", httpCode))
	response.SetHeader("Content-Type", "text/html; charset=utf-8")
	response.SetHeader("X-Content-Type-Options", "nosniff")

	//设置HTTP CODE
	response.SetHttpCode(httpCode)

	//回写HTTP Body
	fmt.Fprintln(response.Writer, ErrorPagesMap[httpCode])
}
