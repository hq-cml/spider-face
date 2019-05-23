package spider

import (
	"net/http"
	"errors"
	"github.com/hq-cml/spider-face/core"
)

type SpiderConfig struct {
	DocRoot       string //web访问目录

	BindAddr      string
	Logger        core.SpiderLogger

	ViewPath      string

	//ReadTimeout   int
	//WriteTimeout  int
	//MaxHeaderByte int

	//HttpErrorHtml map[int]string
}

type Spider struct {
	HttpServer  *http.Server
	MuxHander   http.Handler    //自定义的多路复用器, 替换原生DefaultServerMux, 本质上是一个Handler接口的实现
	Config      *SpiderConfig
}

func NewSpider(sConfig *SpiderConfig) (*Spider, error) { /*{{{*/
	if sConfig.BindAddr == "" {
		return nil, errors.New("server Addr can't be empty...[ip:port]")
	}

	//new Application
	mux := core.NewHandlerMux()

	server := &http.Server{
		Addr: sConfig.BindAddr,
		Handler: mux,
	}

	spd := &Spider{
		Config: sConfig,
		MuxHander: mux,
		HttpServer: server,
	}
	return spd, nil
}


func (spd *Spider) RegisterController(controllerMap map[string]core.SpiderController) {
	mux := spd.MuxHander.(*core.SpiderHandlerMux)
	mux.RegisterController(controllerMap)
}

func (spd *Spider) Run() { /*{{{*/
	//解析模板
	core.CompileTpl(spd.Config.ViewPath)
	//信号处理函数
	//go srv.signalHandle()

	//serverStat = STATE_RUNNING

	//logger.RunLog("[Notice] Server start.")
	//listen loop
	//spd.Serve(srv.listener)
	spd.HttpServer.ListenAndServe()

	//logger.RunLog("[Notice] Waiting for connections to finish...")
	//connWg.Wait()
	//serverStat = STATE_TERMINATE
	//logger.RunLog("[Notice] Server shuttdown.")
	return
}

//TODO
//文件上传 ??
//rewrite ??
//自定制listenner ??
//压缩 ??
//https ??
//热重启??
//MIME??