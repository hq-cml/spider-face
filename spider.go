package spider

import (
	"net/http"
	"errors"
	"github.com/hq-cml/spider-face/core"
	"fmt"
	"github.com/hq-cml/spider-face/utils/helper"
)

type Spider struct {
	HttpServer  *http.Server
	MuxHander   http.Handler    //自定义的多路复用器, 替换原生DefaultServerMux, 本质上是一个Handler接口的实现
	Config      *core.SpiderConfig
}

func NewSpider(sConfig *core.SpiderConfig, controllers map[string]core.SpiderController) (*Spider, error) {
	if sConfig.BindAddr == "" {
		return nil, errors.New("server Addr can't be empty...[ip:port]")
	}

	if sConfig.TplPath == "" {
		sConfig.TplPath = fmt.Sprintf("%s/tpl", helper.GetCurrentDir())
	}

	if sConfig.StaticPath == "" {
		sConfig.StaticPath = fmt.Sprintf("%s/static", helper.GetCurrentDir())
	}

	//new Application
	mux := core.NewHandlerMux(sConfig)

	//注册控制器
	mux.RegisterController(controllers)

	server := &http.Server {
		Addr: sConfig.BindAddr,
		Handler: mux,
	}

	spd := &Spider {
		Config: sConfig,
		MuxHander: mux,
		HttpServer: server,
	}

	return spd, nil
}

func (spd *Spider) Run() {
	//解析模板
	core.CompileTpl(spd.Config.TplPath)
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