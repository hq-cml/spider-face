package core

import "fmt"

type ActionFunc func(rp Roundtrip)

/*
 * Controller接口，规定了Spider中必需的合法的行为
 */
type Controller interface {
	//用户实现
	GetAllRouters() []ControllerRouter

	//Spider提供, 获取Controller的实时状态用于任务分发场景
	Roundtrip
}

type ControllerRouter struct {
	Method   string
	Location string
	Action   string
}

//一个系统默认的controller，用于快捷注册
type DefaultController struct {
	SpiderRoundtrip
	routers    		[]ControllerRouter

	funcMapGet  	map[string]ActionFunc
	funcMapPost  	map[string]ActionFunc
	funcMapPut  	map[string]ActionFunc
	funcMapDelete   map[string]ActionFunc
}

func NewDefaultController() *DefaultController {
	return &DefaultController{
		routers:    	 []ControllerRouter{},
		funcMapGet:		 map[string]ActionFunc{},
		funcMapPost:	 map[string]ActionFunc{},
		funcMapPut: 	 map[string]ActionFunc{},
		funcMapDelete:	 map[string]ActionFunc{},
	}
}

func (def *DefaultController) DefaultGetAction() {
	tmpFunc, ok := def.funcMapGet[def.UrlPath()]
	if ok {
		tmpFunc(def)
	} else {
		fmt.Println("404")
	}
}

func (def *DefaultController) DefaultPostAction() {
	tmpFunc, ok := def.funcMapPost[def.UrlPath()]
	if ok {
		tmpFunc(def)
	} else {
		fmt.Println("404")
	}
}

func (def *DefaultController) DefaultPutAction() {
	tmpFunc, ok := def.funcMapPut[def.UrlPath()]
	if ok {
		tmpFunc(def)
	} else {
		fmt.Println("404")
	}
}

func (def *DefaultController) DefaultDeleteAction() {
	tmpFunc, ok := def.funcMapDelete[def.UrlPath()]
	if ok {
		tmpFunc(def)
	} else {
		fmt.Println("404")
	}
}

func (def *DefaultController) GetAllRouters() []ControllerRouter {
	return def.routers
}

func (mux *HandlerMux) GET(location string , acFunc ActionFunc) {
	defController := mux.DefController.(*DefaultController)
	defController.routers = append(defController.routers, ControllerRouter {
		Method:"GET", Location: location, Action:"DefaultGetAction",
	})
	defController.funcMapGet[location] = acFunc
}

func (mux *HandlerMux) POST(location string , acFunc ActionFunc) {
	defController := mux.DefController.(*DefaultController)
	defController.routers = append(defController.routers, ControllerRouter {
		Method:"POST", Location: location, Action:"DefaultPostAction",
	})
	defController.funcMapPost[location] = acFunc
}

func (mux *HandlerMux) PUT(location string , acFunc ActionFunc) {
	defController := mux.DefController.(*DefaultController)
	defController.routers = append(defController.routers, ControllerRouter {
		Method:"PUT", Location: location, Action:"DefaultPutAction",
	})
	defController.funcMapPut[location] = acFunc
}

func (mux *HandlerMux) DELETE(location string , acFunc ActionFunc) {
	defController := mux.DefController.(*DefaultController)
	defController.routers = append(defController.routers, ControllerRouter {
		Method:"DELETE", Location: location, Action:"DefaultDeleteAction",
	})
	defController.funcMapDelete[location] = acFunc
}