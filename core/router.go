package core

import (
	"strings"
	"errors"
	"reflect"
)

type RouterManager struct {
	logger       SpiderLogger

	ParentMux 	 *SpiderHandlerMux
	RewriteMap	 map[string][]*RouterNode
	RewriteKey	 map[string]string
}

type RouterNode struct {
	urlParts   map[int]map[string]string
	staticNum  int
	method     string
	controller string
	action     string
}

func NewRouterManager(mux *SpiderHandlerMux, logger SpiderLogger) *RouterManager {
	return &RouterManager{
		logger: logger,
		ParentMux: mux,
		RewriteMap: make(map[string][]*RouterNode),
		RewriteKey: make(map[string]string),
	}
}

// Reg routerManger in controller by func RegisterRouter() map[string]interface{}
// support rewrite
func (srt *RouterManager) RegisterRouter(controllerName string, controller SpiderController) error {
	routers := controller.GetRouter()
	if routers == nil {
		return nil
	}

	for _, router := range routers {
		method := router.Method
		pattern := router.Pattern
		action := router.Action

		method = strings.ToUpper(method)
		if old, exist := srt.RewriteKey[method+" "+pattern]; exist {
			return errors.New("Can't register routerManger:\"" + method + " " + pattern + "\",it had been registed in controller:" + old)
		}
		r := createRewriteRouter(method, pattern, controllerName, action)
		if r == nil {
			return errors.New("Can't register routerManger:\"" + method + " " + pattern + "\"")
		}
		srt.RewriteMap[method] = append(srt.RewriteMap[method], r)
		srt.RewriteKey[method+" "+pattern] = controllerName
	}
	return nil
}

// create rewrite routerManger
func createRewriteRouter(method, pattern, controller, action string) *RouterNode {
	if pattern == "" || action == "" {
		return nil
	}
	urlParts := strings.Split(strings.Trim(pattern, "/"), "/")

	router := &RouterNode{
		urlParts:   make(map[int]map[string]string),
		staticNum:  0,
		method:     method,
		controller: controller,
		action:     action,
	}

	for idx, part := range urlParts {
		if part[0:1] == ":" {       //路径参数
			router.urlParts[idx] = map[string]string{
				"name": part[1:],
				"type": "var",
			}
		} else if part[0:1] == "*" { //Pathinfo的参数形式/yera/2019/month/5/day/10
			router.urlParts[idx] = map[string]string{
				"name": "*",
				"type": "",
			}

			break
		} else {
			router.urlParts[idx] = map[string]string{
				"name": part,
				"type": "",
			}
			router.staticNum++
		}
	}
	return router
}

//根据URL和参数，找到对应处理的controller和action
func (srt *RouterManager) MatchRewrite(method, url string) (string, string, map[string]string, error) {
	if _, ok := srt.RewriteMap[method]; ok == false {
		srt.logger.Errf("No support method: %s", method)
		return "", "", nil, errors.New("No match")
	}

	paths := strings.Split(strings.Trim(url, "/"), "/")
	for _, router := range srt.RewriteMap[method] {

		if matchParam, match := srt.matchRouter(router, paths); !match {
			continue
		} else {
			return router.controller, router.action, matchParam, nil
		}
	}

	return "", "", nil, errors.New("No match")
}

func (srt *RouterManager) matchRouter(router *RouterNode, paths []string) (map[string]string, bool) {
	matchParam := map[string]string{}
	var cnt int
	for idx, part := range paths {
		if _, exist := router.urlParts[idx]; !exist {
			return nil, false
		}
		if router.urlParts[idx]["type"] == "" {
			if router.urlParts[idx]["name"] == "*" {
				if idx < len(paths) {
					param := paths[idx:]
					param_num := len(param) / 2
					for i := 0; i < param_num; i++ {
						matchParam[param[i*2]] = param[i*2+1]
					}
				}
				break
			}

			if router.urlParts[idx]["name"] != part {
				return nil, false
			}
			cnt++
		} else if router.urlParts[idx]["type"] == "var" {
			matchParam[router.urlParts[idx]["name"]] = part
		}
	}

	if cnt != router.staticNum {
		return nil, false
	}
	return matchParam, true
}

// Get controller and action name from
// request parame "m"
// eg: demo.index,return demo index
func (this *RouterManager) ParseMethod(method string) (controller_name string, action_name string) {
	method_map := strings.SplitN(method, ".", 2)
	switch len(method_map) {
	case 1:
		controller_name = method_map[0]
	case 2:
		controller_name = method_map[0]
		action_name = method_map[1]
	}
	return controller_name, action_name
}

// create new controller by controller name
func (this *RouterManager) NewController(controllerName string) (reflect.Value, error) {
	//register := GetRegister()

	//m_arr := make([]reflect.Value, 0)
	var newController reflect.Value

	//if register == nil {
	//	//http 500
	//	return controller_instance, errors.New("Server Error : Can't find \"Register\"")
	//}

	controller_type := this.ParentMux.GetController(controllerName)
	if controller_type == nil {
		//http 404
		return newController, errors.New("Warn : Can't find " + controllerName)
	}

	newController = reflect.New(controller_type)
	if false == newController.IsValid() {
		//http 404
		return newController, errors.New("Warn : Can't find " + controllerName)
	}

	return newController, nil
}