package spider

import (
	"net/http"
	"errors"
	"github.com/hq-cml/spider-face/core"
	"fmt"
	"github.com/hq-cml/spider-face/utils/helper"
	"github.com/hq-cml/spider-face/utils/log"
)

type Spider struct {
	HttpServer  *http.Server
	MuxHander   http.Handler    //自定义的多路复用器, 替换原生DefaultServerMux, 本质上是一个Handler接口的实现
	Config      *core.SpiderConfig

	logger      core.SpiderLogger
}

func NewSpider(sConfig *core.SpiderConfig,
	controllers map[string]core.SpiderController, logger core.SpiderLogger) (*Spider, error) {
	if sConfig.BindAddr == "" {
		return nil, errors.New("server Addr can't be empty...[ip:port]")
	}

	if sConfig.TplPath == "" {
		sConfig.TplPath = fmt.Sprintf("%s/tpl", helper.GetCurrentDir())
	}

	if sConfig.StaticPath == "" {
		sConfig.StaticPath = fmt.Sprintf("%s/static", helper.GetCurrentDir())
	}

	if logger == nil {
		if sConfig.LogPath == "" {
			logger = log.DefaultLogger
		} else {
			logger = log.NewLog(sConfig.LogPath, sConfig.LogLevel)
		}
	}

	//new Application
	mux := core.NewHandlerMux(sConfig, logger)

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
		logger : logger,
	}

	spd.logger.Infof("Spider init success!")

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

	spd.logger.Infof("Spider start to run...")
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