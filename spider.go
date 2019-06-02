package spider

import (
	"fmt"
	"net/http"
	"errors"
	"github.com/hq-cml/spider-face/core"
	"github.com/hq-cml/spider-face/utils/helper"
	"github.com/hq-cml/spider-face/utils/log"
)

type Spider struct {
	HttpServer    *http.Server
	MuxHander     http.Handler    //自定义的多路复用器, 替换原生DefaultServerMux, 本质上是一个Handler接口的实现
	Config        *core.SpiderConfig
	logger        core.SpiderLogger
}

//创建spider实例
func NewSpider(sConfig *core.SpiderConfig,
		 logger core.SpiderLogger) (*Spider, error) {

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

	//TODO
	customErrHtml := map[int]string{}
	rewriteRule := map[string]string {
		"/test/rewrite": "/index?name=123",
		"/test/rewrite/(.*)/(.*)": "/index?id=$1&name=$2",
	}

	//创建serverMux
	mux, err := core.NewHandlerMux(sConfig, logger, customErrHtml, rewriteRule)
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

func (spd *Spider) RegisterController(controllers []core.Controller) error {
	//注册控制器
	mux, ok := spd.MuxHander.(*core.HandlerMux)
	if !ok {
		return errors.New("Wrong type of mux!")
	}
	err := mux.RegisterController(controllers)
	if err != nil {
		return err
	}

	return nil
}

func (spd *Spider) Run() {
	//将Default注册到Mux中去
	mux, ok := spd.MuxHander.(*core.HandlerMux)
	if !ok {
		panic("Wrong type of mux!")
	}
	spd.RegisterController([]core.Controller{
		mux.FoolController,
	})

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

func (spd *Spider) GET(location string , acFunc core.ActionFunc) {
	//注册控制器
	mux, ok := spd.MuxHander.(*core.HandlerMux)
	if !ok {
		panic("Wrong type of mux!")
	}
	mux.GET(location, acFunc)
}

func (spd *Spider) POST(location string , acFunc core.ActionFunc) {
	//注册控制器
	mux, ok := spd.MuxHander.(*core.HandlerMux)
	if !ok {
		panic("Wrong type of mux!")
	}
	mux.POST(location, acFunc)
}

func (spd *Spider) PUT(location string , acFunc core.ActionFunc) {
	//注册控制器
	mux, ok := spd.MuxHander.(*core.HandlerMux)
	if !ok {
		panic("Wrong type of mux!")
	}
	mux.PUT(location, acFunc)
}

func (spd *Spider) DELETE(location string , acFunc core.ActionFunc) {
	//注册控制器
	mux, ok := spd.MuxHander.(*core.HandlerMux)
	if !ok {
		panic("Wrong type of mux!")
	}
	mux.DELETE(location, acFunc)
}

//TODO
//文件上传 ??
//rewrite ??
//自定制listenner ??
//压缩 ??
//https ??
//热重启??
//MIME??