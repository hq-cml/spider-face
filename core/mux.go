package core

import (
	"net/http"
	"reflect"
	"fmt"
	"strings"
	"errors"
)

var (
	beforeDispatch = "BeforeDispatch"
	afterDispatch =  "AfterDispatch"
)

//Spider的http多路复用器
//*SpiderHandlerMux实现了http.Handler接口，用来替换掉golang默认的DefaultServerMux
//它是Spider的核心
type HandlerMux struct {
	logger        SpiderLogger
	routerManger  *RouterManager
	controllerMap map[string]reflect.Type
}

//create Application object
func NewHandlerMux(sConfig *SpiderConfig, controllerMap map[string]SpiderController,
		logger SpiderLogger) (*HandlerMux, error) {
	//TODO
	//http_server_config.Root = strings.TrimRight(http_server_config.Root, "/")
	//for err_code, err_file_name := range http_server_config.HttpErrorHtml {
	//	err_html := http_server_config.Root + "/" + strings.TrimLeft(err_file_name, "/")
	//	http_server_config.HttpErrorHtml[err_code] = err_html
	//}

	//用外层用户定制的conf初始化全局配置
	GlobalConf = sConfig

	//生成mux
	mux := &HandlerMux{
		logger: logger,
		controllerMap: map[string]reflect.Type{},
	}

	//init mime
	//initMime()

	//创建路由
	mux.routerManger = NewRouterManager(mux, logger)

	//注册控制器
	err := mux.RegisterController(controllerMap)
	if err != nil {
		return nil, err
	}

	mux.logger.Info("Server mux done~")
	return mux, nil
}

//注册控制器
//TODO 应该有个默认的Controller
func (mux *HandlerMux) RegisterController(controllerMap map[string]SpiderController) error {
	for name, controller := range controllerMap {
		//验重
		if _, exist := mux.controllerMap[name]; exist {
			mux.logger.Errf("Conflicting controller: %v", name)
			return fmt.Errorf("Controller %q is existed!", name)
		}

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

//通过反射包和controller的名字，还原Controller的reflect.Value
func (mux *HandlerMux) ValueOfController(controllerName string) (reflect.Value, error) {
	var valueOfController reflect.Value

	typeOfController, exist := mux.controllerMap[controllerName];
	if !exist {
		//http 404
		return valueOfController, errors.New("Warn : Can't find " + controllerName)
	}

	valueOfController = reflect.New(typeOfController)
	if !valueOfController.IsValid() {
		//http 404
		return valueOfController, errors.New("Warn : Can't find " + controllerName)
	}

	return valueOfController, nil
}

//实现http.Handler接口
func (mux *HandlerMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//start_time := time.Now()

	//TODO rewrite
	//if r.URL.Path != "/" {
	//	matchRewrite(r)
	//}

	mux.DispatchHandler(w, r)

	//end_time := time.Now()
	//
	//request_time := float64(end_time.UnixNano()-start_time.UnixNano()) / 1000000000
	//
	//log_format := "%s - [%s] %s %s %s %s %.5f \"%s\"" //ip - [time] Method uri scheme status request_time agent

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

func (mux *HandlerMux) DispatchHandler(w http.ResponseWriter, r *http.Request) {
	routerManager := mux.routerManger
	//init request & reponse
	request := NewRequest(r)
	response := NewResponse(w, request)

	var controllerName string
	var actionName string
	var pathParam map[string]string
	var ok error

	//去除urlPaht，交给路由管理利器来分析，得到controller和action
	urlPath := strings.TrimRight(request.UrlPath(), "/")
	fmt.Println("REQ URL PATH: ", urlPath)
	if urlPath != "" { //有url
		controllerName, actionName, pathParam, ok = routerManager.AnalysePath(r.Method, urlPath)
		if ok != nil {
			//分析失败，可能是没有找到合适的处理器，也可能是一个静态文件
			OutputStaticFile(response, request, urlPath)
			return
		}

		actionName = strings.TrimSuffix(actionName, ACTION_SUFFIX)
	} else  { //首页
		controllerName = DEFAULT_CONTROLLER
		actionName = DEFAULT_ACTION
	}

	request.rewriteParams = pathParam
	request.SetController(controllerName)
	request.SetAction(actionName)

	valueOfController, err := mux.ValueOfController(controllerName)
	if err != nil {
		OutErrorHtml(response, request, http.StatusNotFound)
		fmt.Println(err)
		return
	}

	controllerHandler := valueOfController.MethodByName(actionName + ACTION_SUFFIX)
	if controllerHandler.IsValid() == false {
		OutErrorHtml(response, request, http.StatusNotFound)
		return
	}

	initParams := make([]reflect.Value, 2)
	initParams[0] = reflect.ValueOf(request)
	initParams[1] = reflect.ValueOf(response)

	initIandler := valueOfController.MethodByName("Init")
	if initIandler.IsValid() == false {
		//logger.ErrorLog("Can't find Method of \"Init\" in controller " + controller_name)
		//OutErrorHtml(response, request, http.StatusInternalServerError)
		panic("A")
		return
	}

	handlers := make([]reflect.Value, 0)
	if beforeHandler := valueOfController.MethodByName(beforeDispatch); beforeHandler.IsValid() == true {
		handlers = append(handlers, beforeHandler)
	}

	handlers = append(handlers, controllerHandler)
	if afterHandler := valueOfController.MethodByName(afterDispatch); afterHandler.IsValid() == true {
		handlers = append(handlers, afterHandler)
	}

	//执行 Init()
	initResult := initIandler.Call(initParams)

	if reflect.Indirect(initResult[0]).Bool() == false {
		//logger.ErrorLog("Method of \"Init\" in controller " + controller_name + " return false")
		OutErrorHtml(response, request, http.StatusInternalServerError)
		return
	}

	requestParams := make([]reflect.Value, 0)
	//Run : Init -> before_dispatch -> controller_handler -> after_dispatch
	for _, v := range handlers {
		v.Call(requestParams)
	}

	response.Header("Connection", request.Header("Connection"))
}