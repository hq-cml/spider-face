package core

import (
	"reflect"
	"strings"
	"fmt"
	"net/http"
)

type Dispatcher struct {
	beforeDispatch string
	afterDispatch  string
}

func NewDispatcher() *Dispatcher {
	return &Dispatcher{
		beforeDispatch: "BeforeDispatch",
		afterDispatch:  "AfterDispatch",
	}
}

func (this *Dispatcher) DispatchHandler(srt *SpiderRouter, w http.ResponseWriter, r *http.Request) {
	//init request
	request := NewRequest(r)
	response := NewResponse(w, request)

	var controllerName string
	var actionName string
	var matchParam map[string]string
	var ok error

	urlPath := strings.TrimRight(request.UrlPath(), "/")
	fmt.Println("r.URL PATH: ", urlPath)
	if urlPath != "" { //有url
		controllerName, actionName, matchParam, ok = srt.MatchRewrite(r.Method, urlPath)
		if ok != nil {
			fmt.Println("A---------------------", urlPath)
			OutputStaticFile(response, request, urlPath)
			return
		}

		if matchParam != nil && len(matchParam)>0{
			request.rewriteParams = matchParam
		}
		actionName = strings.TrimSuffix(actionName, ACTION_SUFFIX)
	} else if urlPath == "" && request.Param(HTTP_METHOD_PARAM_NAME) == "" { //首页
		controllerName = DEFAULT_CONTROLLER
		actionName = DEFAULT_ACTION
	} else if request.Param(HTTP_METHOD_PARAM_NAME) != "" {
		controllerName, actionName = srt.ParseMethod(request.Param(HTTP_METHOD_PARAM_NAME))
		actionName = strings.Title(strings.ToLower(actionName))
	}

	request.SetController(controllerName)
	request.SetAction(actionName)

	controller, err := srt.NewController(controllerName)
	if err != nil {
		OutErrorHtml(response, request, http.StatusNotFound)
		fmt.Println(err)
		return
	}

	controllerHandler := controller.MethodByName(actionName + ACTION_SUFFIX)
	if controllerHandler.IsValid() == false {
		OutErrorHtml(response, request, http.StatusNotFound)
		fmt.Println(err)
		return
	}

	initParams := make([]reflect.Value, 2)
	initParams[0] = reflect.ValueOf(request)
	initParams[1] = reflect.ValueOf(response)

	initIandler := controller.MethodByName("Init")
	if initIandler.IsValid() == false {
		//logger.ErrorLog("Can't find method of \"Init\" in controller " + controller_name)
		//OutErrorHtml(response, request, http.StatusInternalServerError)
		panic("A")
		return
	}

	handlers := make([]reflect.Value, 0)
	if beforeHandler := controller.MethodByName(this.beforeDispatch); beforeHandler.IsValid() == true {
		handlers = append(handlers, beforeHandler)
	}

	handlers = append(handlers, controllerHandler)
	if afterHandler := controller.MethodByName(this.afterDispatch); afterHandler.IsValid() == true {
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
