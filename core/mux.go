package core

import (
	"net/http"
	"reflect"
	"fmt"
)

//Spider的http多路复用器
//*SpiderHandlerMux实现了http.Handler接口，用来替换掉golang默认的DefaultServerMux
//它是Spider的核心
type SpiderHandlerMux struct {
	logger        SpiderLogger
	dispatcher    *Dispatcher
	routerManger  *RouterManager
	controllerMap map[string]reflect.Type
}

//create Application object
func NewHandlerMux(sConfig *SpiderConfig, controllers map[string]SpiderController,
		logger SpiderLogger) (*SpiderHandlerMux, error) {
	//TODO
	//http_server_config.Root = strings.TrimRight(http_server_config.Root, "/")
	//for err_code, err_file_name := range http_server_config.HttpErrorHtml {
	//	err_html := http_server_config.Root + "/" + strings.TrimLeft(err_file_name, "/")
	//	http_server_config.HttpErrorHtml[err_code] = err_html
	//}

	//用外层用户定制的conf初始化全局配置
	GlobalConf = sConfig

	//生成mux
	mux := &SpiderHandlerMux{
		logger: logger,
		controllerMap: map[string]reflect.Type{},
	}

	//init mime
	//initMime()

	//init dispatcher
	//TODO 待优化
	mux.dispatcher = NewDispatcher()

	//创建路由
	mux.routerManger = NewRouterManager(mux, logger)

	//注册控制器
	err := mux.RegisterController(controllers)
	if err != nil {
		return nil, err
	}

	mux.logger.Info("Server mux done~")
	return mux, nil
}

//注册控制器
//TODO 应该有个默认的Controller
func (mux *SpiderHandlerMux) RegisterController(controllerMap map[string]SpiderController) error {
	for name, controller := range controllerMap {
		//验重
		if _, exist := mux.controllerMap[name]; exist {
			mux.logger.Errf("Conflicting controller: %v", name)
			return fmt.Errorf("Controller %q is existed!", name)
		}

		//var i interface{}
		//var a interface{}
		//a = 10
		//i = &a
		//
		//fmt.Println(reflect.ValueOf(a))
		//fmt.Println(reflect.ValueOf(i))
		//
		//fmt.Println(reflect.Indirect(reflect.ValueOf(a)))
		//fmt.Println(reflect.Indirect(reflect.ValueOf(i)))

		//获取controller的reflect.Value值
		//reflect.Indirect保证即便是指针也能拿到实际的指向值
		controllerValue := reflect.Indirect(reflect.ValueOf(controller))
		mux.controllerMap[name] = controllerValue.Type()

		//将各controller的路由注册上来
		err := mux.routerManger.RegisterRouter(name, controller)
		if err != nil {
			mux.logger.Errf("RegController error :%v", err)
			return err
		}
	}

	return nil
}

func (this *SpiderHandlerMux) GetController(controllerName string) reflect.Type {
	if c, ok := this.controllerMap[controllerName]; ok == false {
		return nil
	} else {
		return c.(reflect.Type)
	}
}

//实现http.Handler接口
func (mux *SpiderHandlerMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//start_time := time.Now()

	//TODO rewrite
	//if r.URL.Path != "/" {
	//	matchRewrite(r)
	//}

	mux.dispatcher.DispatchHandler(mux.routerManger, w, r)

	//end_time := time.Now()
	//
	//request_time := float64(end_time.UnixNano()-start_time.UnixNano()) / 1000000000
	//
	//log_format := "%s - [%s] %s %s %s %s %.5f \"%s\"" //ip - [time] method uri scheme status request_time agent

	//access_log := fmt.Sprintf(log_format,
	//	mux.Isset(r.RemoteAddr),
	//	Date("Y/m/d H:i:s", start_time),
	//	mux.Isset(r.Method),
	//	mux.Isset(r.URL.RequestURI()),
	//	mux.Isset(r.Proto),
	//	mux.Isset(w.Header().Get("Status")),
	//	request_time,
	//	mux.Isset(r.Header.Get("User-Agent")),
	//)
	//logger.AccessLog(access_log)
}