package core

import (
	"net/http"
	"reflect"
)

//Spider的http多路复用器
type SpiderHandlerMux struct {
	Logger        SpiderLogger
	dispatcher    *Dispatcher
	router        *SpiderRouter
	controllerMap     map[string]interface{}
}

//create Application object
func NewHandlerMux() *SpiderHandlerMux {
	//TODO
	//http_server_config.Root = strings.TrimRight(http_server_config.Root, "/")
	//for err_code, err_file_name := range http_server_config.HttpErrorHtml {
	//	err_html := http_server_config.Root + "/" + strings.TrimLeft(err_file_name, "/")
	//	http_server_config.HttpErrorHtml[err_code] = err_html
	//}

	mux := &SpiderHandlerMux{
		controllerMap: map[string]interface{}{},
	}

	//init mime
	//initMime()

	//init dispatcher
	mux.dispatcher = NewDispatcher()

	//init router
	mux.router = NewRouter(mux)

	return mux
}

func (mux *SpiderHandlerMux) RegisterController(controllerMap map[string]SpiderController) {
	for name, controller := range controllerMap {
		if _, exist := mux.controllerMap[name]; exist {
			//logger.RunLog("[Error] conflicting controller name:" + controller_name)
			//return fmt.Errorf("%q is existed!", name)
			continue;
		}

		controllerValue := reflect.Indirect(reflect.ValueOf(controller))
		mux.controllerMap[name] = controllerValue.Type()

		err := mux.router.RegRouter(name, controller)
		if err != nil {
			//logger.RunLog(fmt.Sprintf("[Error] RegController error :%v", err))
			//os.Exit(0)
			panic(err)
		}
	}
}

func (this *SpiderHandlerMux) GetController(controller_name string) reflect.Type {
	if c, ok := this.controllerMap[controller_name]; ok == false {
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

	mux.dispatcher.DispatchHandler(mux.router, w, r)

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