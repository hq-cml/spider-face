package core

import "fmt"

//一个系统默认的controller，用于快捷注册
type DefaultController struct {
	RuntimeController
	routers    []ControllerRouter
	funcMapGet  map[string]ActionFunc
}

func (def *DefaultController) DefaultGetAction() {
	//def.Echo("hello world!")
	fmt.Println("A-----------", def.UrlPath())
	fmt.Println("B-----------", def.GetUri())

	tmpFunc, ok := def.funcMapGet["/hello"]
	if ok {
		tmpFunc(def)
	} else {
		fmt.Println("?????")
	}

}

func (def *DefaultController) GetRouter() []ControllerRouter {
	return def.routers
}

type ActionFunc func(c Controller)

func (mux *HandlerMux) GET(location string , acFunc ActionFunc) {
	defController := mux.DefController.(*DefaultController)
	defController.routers = append(defController.routers, ControllerRouter {
		Method:"GET", Location: location, Action:"DefaultGetAction",
	})
	fmt.Println("X-------------", location)
	defController.funcMapGet[location] = acFunc
}