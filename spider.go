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

//创建spider实例
func NewSpider(sConfig *core.SpiderConfig,
		controllerMap map[string]core.Controller, logger core.SpiderLogger) (*Spider, error) {

	if sConfig.BindAddr == "" {
		return nil, errors.New("server Addr can't be empty...[ip:port]")
	}

	//如果没有模板路径和静态文件路径，则用默认的
	if sConfig.TplPath == "" {
		sConfig.TplPath = fmt.Sprintf("%s/tpl", helper.GetCurrentDir())
	}
	if sConfig.StaticPath == "" {
		sConfig.StaticPath = fmt.Sprintf("%s/static", helper.GetCurrentDir())
	}

	//如果用户没有给定logger，则使用默认的logger
	if logger == nil {
		if sConfig.LogPath == "" {
			logger = log.DefaultLogger
		} else {
			logger = log.NewLog(sConfig.LogPath, sConfig.LogLevel)
		}
	}

	//创建serverMux
	mux, err := core.NewHandlerMux(sConfig, controllerMap, logger)
	if err != nil {
		return nil, err
	}

	//替换掉golang自带的handler
	server := &http.Server {
		Addr: sConfig.BindAddr,
		Handler: mux,
	}

	//创建spider实例
	spd := &Spider {
		Config: sConfig,
		MuxHander: mux,
		HttpServer: server,
		logger : logger,
	}

	//初始化解析视图模板
	err = core.InitViewTemplate(spd.Config.TplPath, spd.logger)
	if err != nil {
		return nil, err
	}

	spd.logger.Info("Spider init success!")

	return spd, nil
}

func (spd *Spider) Run() {
	//信号处理函数
	//go srv.signalHandle()

	//serverStat = STATE_RUNNING

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