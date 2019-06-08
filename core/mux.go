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
	"time"
)

//TODO 这个扩展成通用的Hook
var (
	beforeDispatch = "BeforeDispatchHook"
	afterDispatch =  "AfterDispatchHook"
)

//多路复用器，用来替换掉golang默认的DefaultServerMux
type HandlerMux struct {
	logger           SpiderLogger
	rewriter         *Rewriter               //地址重写器，地址重写工作在路由之前完成
    customErrHtml    map[int]string
	routerManger     *RouterManager          //路由管理器，负责实际的请求路由到对应的Controller/Action
	controllerMap    map[string]reflect.Type //所有controller的动态类型，用于controller实例还原

	SpeedyController *SpeedyController       //Spider自带的默认Controller，用于快捷注册使用
}

//create Application object
func NewHandlerMux(sConfig *SpiderConfig, logger SpiderLogger,
		customErrHtmls map[int]string, rewriteRule map[string]string) (*HandlerMux) {

	//初始化用户自定义的错误页面,如果有
	customErrHtml := map[int]string{}
	for code, errHtml := range customErrHtmls {
		ht := sConfig.StaticPath + "/" + strings.TrimLeft(errHtml, "/")
		customErrHtml[code] = ht
	}

	//用外层用户定制的conf初始化全局配置
	GlobalConf = sConfig

	//生成mux
	mux := &HandlerMux {
		logger:           logger,
		customErrHtml:    customErrHtml,
		controllerMap:    map[string]reflect.Type{},
		SpeedyController: NewSpeedyController(),
	}

	//生成rewriter
	mux.rewriter = NewRewriter(logger)

	//注册重写规则
	mux.rewriter.RegisterRewriteRule(rewriteRule)

	//init mime
	//initMime()

	//创建路由
	mux.routerManger = NewRouterManager(logger)

	mux.logger.Info("ServerMux Init ~")
	return mux
}

//注册控制器
func (mux *HandlerMux) RegisterController(controllers []Controller) error {
	for _, controller := range controllers {
		//获取controller的名字
		typ := reflect.Indirect(reflect.ValueOf(controller)).Type()
		if strings.Index(typ.Name(), CONTROLLER_SUFFIX) == -1 {
			return errors.New("Invalid Controller Name! Must End with 'Controller'")
		}
		name := strings.TrimSuffix(typ.Name(), CONTROLLER_SUFFIX)
		//名字验重
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
		mux.logger.Infof("Register Controller : %v", typ.Name())
	}

	return nil
}

//实现http.Handler接口
func (mux *HandlerMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	//尝试匹配rewrite规则（如果有）
	if r.URL.Path != "/" {
		mux.rewriter.TryMatchRewrite(r)
	}

	//分派处理
	mux.DispatchHandler(w, r)

	//收尾日志
	endTime := time.Now()
	costDuration := float64(endTime.UnixNano() - startTime.UnixNano()) / 1000000 //毫秒
	format := "AccessLog: Method=[%s]; Uri=[%s]; ClientIp=[%s]; StartTime=[%s]; CostTime=[%f] ms; Status=[%s]"
	mux.logger.Infof(format, r.Method, r.URL.RequestURI(),
		r.RemoteAddr, Date("Y-m-d H:i:s", startTime), costDuration,
		w.Header().Get("Status")/*, r.Header.Get("User-Agent")*/)
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
	urlPath := ""
	if request.UrlPath() != "/" {
		urlPath = strings.TrimRight(request.UrlPath(), "/")
	} else {
		urlPath = request.UrlPath()
	}
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

	if controllerName == SPEEDY_CONTROLLER_NAME {
		mux.handleSpeedyController(request, response, controllerName, actionName)
	} else {
		mux.handleNormalController(request, response, controllerName, actionName)
	}

	response.SetHeader("Connection", request.GetHeader("Connection"))
}

func (mux *HandlerMux) handleSpeedyController(request *Request, response *Response,
	controllerName, actionName string) {

	//还原创建出一个Controller实时实例
	spdController := SpeedyController{
		funcMapGet: mux.SpeedyController.funcMapGet,
		funcMapPost: mux.SpeedyController.funcMapPost,
		funcMapPut: mux.SpeedyController.funcMapPut,
		funcMapDelete: mux.SpeedyController.funcMapDelete,
	}
	//还原创建Controller实例的roundtrip
	spdController.GetRoundTrip().initRoundtrip(request, response, controllerName, actionName, mux.logger)

	switch request.GetMethod() {
	case http.MethodGet:
		spdController.SpeedyGetAction()
	case http.MethodPost:
		spdController.SpeedyPostAction()
	case http.MethodPut:
		spdController.SpeedyPutAction()
	case http.MethodDelete:
		spdController.SpeedyDeleteAction()
	default:
		spdController.SpeedyGetAction()
	}

	//TODO beforeAction
	//TODO afterAction
}

func (mux *HandlerMux) handleNormalController(request *Request, response *Response,
		controllerName, actionName string) {
	valueOfController, err := mux.ValueOfController(controllerName)
	if err != nil {
		OutputErrorHtml(response, request, http.StatusNotFound, mux.customErrHtml)
		fmt.Println(err)
		return
	}

	//执行controller的Init()
	c, b := valueOfController.Interface().(Controller)
	if !b {
		panic("Oh my god")
	}
	rp := c.GetRoundTrip()
	rp.initRoundtrip(request, response, controllerName, actionName, mux.logger)

	//执行Action（包括前后的Hook，如果有）
	actions := make([]reflect.Value, 0)
	controllerAction := valueOfController.MethodByName(actionName + ACTION_SUFFIX)
	if controllerAction.IsValid() == false {
		OutputErrorHtml(response, request, http.StatusNotFound, mux.customErrHtml)
		return
	}
	if beforeAction := valueOfController.MethodByName(beforeDispatch); beforeAction.IsValid() == true {
		actions = append(actions, beforeAction)
	}
	actions = append(actions, controllerAction)
	if afterAction := valueOfController.MethodByName(afterDispatch); afterAction.IsValid() == true {
		actions = append(actions, afterAction)
	}

	//roundTrip就是Action函数的参数
	requestParams := []reflect.Value{
		reflect.ValueOf(rp),
	}
	//Run : before_dispatch -> controller_handler -> after_dispatch
	for _, action := range actions {
		action.Call(requestParams)
	}
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

