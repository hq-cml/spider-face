package core

import (
	"reflect"
	"strings"
	"fmt"
	"net/http"
)

type Dispatcher struct {
	before_dispatch string
	after_dispatch  string
}

func NewDispatcher() *Dispatcher {
	return &Dispatcher{
		before_dispatch: "BeforeDispatch",
		after_dispatch:  "AfterDispatch",
	}
}

func (this *Dispatcher) DispatchHandler(srt *SpiderRouter, w http.ResponseWriter, r *http.Request) {
	//init request
	request := NewRequest(r)
	//router := srt.GetRouter()
	response := NewResponse(w, request)

	var controller_name string
	var actionName string
	var matchParam map[string]string
	var ok error

	_ = matchParam //TODO
	//w.Header("Status", fmt.Sprintf("%d", http.StatusOK))
	w.WriteHeader(http.StatusOK)

	url := strings.TrimRight(request.Url(), "/")
	fmt.Println("r.URL: ", url)
	if url != "" { //有url
		controller_name, actionName, matchParam, ok = srt.MatchRewrite(url, r.Method)
		if ok != nil {
			OutputStaticFile(response, request, url)
			return
		}

		//TODO
		//if match_param != nil {
		//	request.rewrite_params = match_param
		//}
		actionName = strings.TrimSuffix(actionName, ACTION_SUFFIX)
	} else if url == "" && request.Param(HTTP_METHOD_PARAM_NAME) == "" { //首页
		controller_name = DEFAULT_CONTROLLER
		actionName = DEFAULT_ACTION
	} else if request.Param(HTTP_METHOD_PARAM_NAME) != "" {
		controller_name, actionName = srt.ParseMethod(request.Param(HTTP_METHOD_PARAM_NAME))
		actionName = strings.Title(strings.ToLower(actionName))
	}

	request.SetController(controller_name)
	request.SetAction(actionName)

	controller, err := srt.NewController(controller_name)
	if err != nil {
		//OutErrorHtml(response, request, http.StatusNotFound)
		fmt.Println(err)
		return
	}

	controllerHandler := controller.MethodByName(actionName + ACTION_SUFFIX)
	if controllerHandler.IsValid() == false {
		//OutErrorHtml(response, request, http.StatusNotFound)
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

	request_handlers := make([]reflect.Value, 0)
	//TODO
	//if before_handler := controller.MethodByName(this.before_dispatch); before_handler.IsValid() == true {
	//	request_handlers = append(request_handlers, before_handler)
	//}

	request_handlers = append(request_handlers, controllerHandler)
	//TODO
	//if after_handler := controller.MethodByName(this.after_dispatch); after_handler.IsValid() == true {
	//	request_handlers = append(request_handlers, after_handler)
	//}

	//执行 Init()
	init_result := initIandler.Call(initParams)

	if reflect.Indirect(init_result[0]).Bool() == false {
		//logger.ErrorLog("Method of \"Init\" in controller " + controller_name + " return false")
		//OutErrorHtml(response, request, http.StatusInternalServerError)
		return
	}

	requestParams := make([]reflect.Value, 0)
	//Run : Init -> before_dispatch -> controller_handler -> after_dispatch
	for _, v := range request_handlers {
		v.Call(requestParams)
	}

	response.Header("Connection", request.Header("Connection"))
}
