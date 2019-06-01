package core

/*
 * 实现自己的多路复用器
 * 它会替换掉golang默认的DefaultServerMux，是Spider的核心
 */

import (
	"net/http"
	"reflect"
	"fmt"
	"strings"
	"errors"
)

//TODO 这个扩展成通用的Hook
var (
	beforeDispatch = "BeforeDispatchHook"
	afterDispatch =  "AfterDispatchHook"
)

//多路复用器，用来替换掉golang默认的DefaultServerMux
type HandlerMux struct {
	logger        SpiderLogger
	rewriter  	  *Rewriter
    customErrHtml map[int]string
	routerManger  *RouterManager
	controllerMap map[string]reflect.Type
}

//create Application object
func NewHandlerMux(sConfig *SpiderConfig, controllers []Controller,
		logger SpiderLogger, customErrHtmls map[int]string,
		rewriteRule map[string]string) (*HandlerMux, error) {

	//初始化用户自定义的错误页面,如果有
	customErrHtml := map[int]string{}
	for code, errHtml := range customErrHtmls {
		ht := sConfig.StaticPath + "/" + strings.TrimLeft(errHtml, "/")
		customErrHtml[code] = ht
	}

	//用外层用户定制的conf初始化全局配置
	GlobalConf = sConfig

	//生成mux
	mux := &HandlerMux{
		logger: logger,
		customErrHtml:customErrHtml,
		controllerMap: map[string]reflect.Type{},
	}

	//生成rewriter
	mux.rewriter = NewRewriter(logger)

	//注册重写规则
	mux.rewriter.RegisterRewriteRule(rewriteRule)

	//init mime
	//initMime()

	//创建路由
	mux.routerManger = NewRouterManager(logger)

	//注册控制器
	err := mux.RegisterController(controllers)
	if err != nil {
		return nil, err
	}

	mux.logger.Info("ServerMux Init ~")
	return mux, nil
}

//注册控制器
//TODO 应该有个默认的Controller
func (mux *HandlerMux) RegisterController(controllers []Controller) error {
	for _, controller := range controllers {
		//获取controller的名字
		typ := reflect.Indirect(reflect.ValueOf(controller)).Type()
		if strings.Index(typ.Name(), "Controller") == -1 {
			return errors.New("Invalid Controller Name! Must End with 'Controller'")
		}
		//name := strings.ToLower(strings.TrimSuffix(typ.Name(), "Controller"))
		name := strings.TrimSuffix(typ.Name(), "Controller")

		fmt.Println("A-----------", name)
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

	fmt.Println("B------------", controllerName)
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

	//尝试匹配rewrite规则
	if r.URL.Path != "/" {
		mux.rewriter.TryMatchRewrite(r)
	}

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
	response := NewResponse(w)

	var controllerName string
	var actionName string
	var pathParam map[string]string
	var ok error

	//去除urlPaht，交给路由管理利器来分析，得到controller和action
	urlPath := strings.TrimRight(request.UrlPath(), "/")
	fmt.Println("REQ URL PATH: ", r.Method, urlPath)
	if urlPath != "" { //有url
		controllerName, actionName, pathParam, ok = routerManager.AnalysePath(r.Method, urlPath)
		if ok != nil {
			//分析失败，可能是没有找到合适的处理器，也可能是一个静态文件
			OutputStaticFile(response, request, urlPath, mux.customErrHtml)
			return
		}

		actionName = strings.TrimSuffix(actionName, ACTION_SUFFIX)
	} else  { //首页
		controllerName = DEFAULT_CONTROLLER
		actionName = DEFAULT_ACTION
	}

	request.pathParams = pathParam

	valueOfController, err := mux.ValueOfController(controllerName)
	if err != nil {
		OutErrorHtml(response, request, http.StatusNotFound, mux.customErrHtml)
		fmt.Println(err)
		return
	}

	//执行controller的Init()
	c, b := valueOfController.Interface().(Controller)
	if !b {
		panic("Oh my god")
	}
	c.Init(request, response, mux.logger)
	c.SetName(controllerName)
	c.SetAction(actionName)

	//执行Action（包括前后的Hook，如果有）
	actions := make([]reflect.Value, 0)
	controllerAction := valueOfController.MethodByName(actionName + ACTION_SUFFIX)
	if controllerAction.IsValid() == false {
		OutErrorHtml(response, request, http.StatusNotFound, mux.customErrHtml)
		return
	}
	if beforeAction := valueOfController.MethodByName(beforeDispatch); beforeAction.IsValid() == true {
		actions = append(actions, beforeAction)
	}
	actions = append(actions, controllerAction)
	if afterAction := valueOfController.MethodByName(afterDispatch); afterAction.IsValid() == true {
		actions = append(actions, afterAction)
	}

	requestParams := make([]reflect.Value, 0)
	//Run : before_dispatch -> controller_handler -> after_dispatch
	for _, action := range actions {
		action.Call(requestParams)
	}

	response.SetHeader("Connection", request.GetHeader("Connection"))
}