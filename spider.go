package spider

import (
	"fmt"
	"time"
	"net/http"
	"github.com/hq-cml/spider-face/core"
	"github.com/hq-cml/spider-face/utils/helper"
	"github.com/hq-cml/spider-face/utils/log"
	"os"
	"os/signal"
	"syscall"
)

type Spider struct {
	HttpServer    *http.Server
	MuxHander     *core.HandlerMux    //自定义的多路复用器, 替换原生DefaultServerMux, 本质上是一个Handler接口的实现
	Config        *core.SpiderConfig
	logger        core.SpiderLogger
}

//创建spider实例
func NewSpider(sConfig *core.SpiderConfig, logger core.SpiderLogger) (*Spider) {
	//默认的config实例
	if sConfig == nil {
		sConfig = &core.SpiderConfig{}
	}

	//如果用户没有给定logger，则使用默认的logger
	if logger == nil {
		if sConfig.LogPath == "" {
			logger = log.DefaultLogger
		} else {
			logger = log.NewLog(sConfig.LogPath, sConfig.LogLevel)
		}
	}

	//如果没有给定地址，那么就用默认的
	if sConfig.BindAddr == "" {
		logger.Info("No Addr. Use Default :9529")
		sConfig.BindAddr = ":9529"
	}

	//如果没有模板路径和静态文件路径，则用默认的
	if sConfig.TplPath == "" {
		sConfig.TplPath = fmt.Sprintf("%s/tpl", helper.GetCurrentDir())
	}
	if sConfig.StaticPath == "" {
		sConfig.StaticPath = fmt.Sprintf("%s/static", helper.GetCurrentDir())
	}

	//用户自定义的错误页面和重写规则
	if sConfig.CustomHttpErrorHtml == nil {
		sConfig.CustomHttpErrorHtml = map[int]string{}
	}
	if sConfig.CustomRewriteRule != nil {
		//"/test/rewrite" => "/index?name=123",
		//"/test/rewrite/(.*)/(.*)" => "/index?id=$1&name=$2",
		sConfig.CustomRewriteRule = map[string]string{}
	}

	//用户自定义的timeout和maxheader
	if sConfig.ReadTimeout <= 0 {
		sConfig.ReadTimeout = 30
	}
	if sConfig.WriteTimeout <= 0 {
		sConfig.WriteTimeout = 30
	}
	if sConfig.MaxHeaderByte <= 0 {
		sConfig.MaxHeaderByte = 1 << 20
	}

	//创建serverMux
	mux := core.NewHandlerMux(sConfig, logger,
		sConfig.CustomHttpErrorHtml, sConfig.CustomRewriteRule)

	//创建http.server实例，替换掉golang自带的handler
	server := &http.Server {
		Addr: sConfig.BindAddr,
		Handler: mux,

		ReadTimeout:    time.Duration(sConfig.ReadTimeout * int64(time.Second)),
		WriteTimeout:   time.Duration(sConfig.WriteTimeout * int64(time.Second)),
		MaxHeaderBytes: 1 << 20,
	}

	//创建spider实例
	spd := &Spider {
		Config: sConfig,
		MuxHander: mux,
		HttpServer: server,
		logger : logger,
	}

	//初始化解析视图模板文件
	core.InitViewTemplate(spd.Config.TplPath, spd.logger)

	spd.logger.Debugf("TplPath: %v", sConfig.TplPath)
	spd.logger.Debugf("StaticPath: %v", sConfig.StaticPath)
	spd.logger.Info("Spider init success!")
	return spd
}

//注册controller
func (spd *Spider) RegisterController(controllers []core.Controller) error {
	//注册控制器
	err := spd.MuxHander.RegisterController(controllers)
	if err != nil {
		return err
	}

	return nil
}

//Run
func (spd *Spider) Run() {
	//将Default注册到Mux中去
	err := spd.RegisterController([]core.Controller{
		spd.MuxHander.SpeedyController,
	})
	if err != nil {
		panic(err)
	}

	//异步信号处理
	go spd.signalHandle()

	//listen loop
	spd.logger.Infof("Spider start to run...")
	spd.HttpServer.ListenAndServe()

	spd.logger.Infof("The Http Server is Closed!")
	spd.logger.Info("Bye Bye~")
	return
}

//异步信号处理
func (spd *Spider) signalHandle() {
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGHUP, syscall.SIGTERM, syscall.SIGQUIT)

	for {
		sig := <-ch
		switch sig {
		case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
			//TODO 强行关闭不优雅，后续改成Shutdown
			spd.HttpServer.Close()
			spd.logger.Infof("Spider Recv Signal: %v", sig)
		default:
			spd.logger.Infof("Unknown signal: %v", sig)
		}
	}
}

//快捷注册，向SpeedController中注入路由规则
func (spd *Spider) GET(location string , acFunc core.ActionFunc) {
	spd.MuxHander.SpeedyController.GET(location, acFunc)
}

func (spd *Spider) POST(location string , acFunc core.ActionFunc) {
	spd.MuxHander.SpeedyController.POST(location, acFunc)
}

func (spd *Spider) PUT(location string , acFunc core.ActionFunc) {
	spd.MuxHander.SpeedyController.PUT(location, acFunc)
}

func (spd *Spider) DELETE(location string , acFunc core.ActionFunc) {
	spd.MuxHander.SpeedyController.DELETE(location, acFunc)
}

//TODO
//文件上传 ??
//rewrite ??
//自定制listenner ??
//压缩 ??
//https ??
//热重启??
//优雅退出Shutdown + context
//MIME??