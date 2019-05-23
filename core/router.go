package core

import (
	"strings"
	"errors"
	"reflect"
)

type SpiderRouter struct {
	ParentMux *SpiderHandlerMux
	RewriteMap map[string][]*RewriteRouter
	RewriteKey map[string]string
}

type RewriteRouter struct {
	urlParts   map[int]map[string]string
	staticNum  int
	method     string
	controller string
	action     string
}

func NewRouter(mux *SpiderHandlerMux) *SpiderRouter {
	return &SpiderRouter{
		ParentMux: mux,
		RewriteMap: make(map[string][]*RewriteRouter),
		RewriteKey: make(map[string]string),
	}
}

// Reg router in controller by func RegRouter() map[string]interface{}
// support rewrite
func (srt *SpiderRouter)RegRouter(controllerName string, controller SpiderController ) error {
	urlMap := controller.GetRouter()
	if urlMap == nil {
		return nil
	}

	for pattern, tb := range urlMap {
		switch tb.(type) {
		case string:
			if old, exist := srt.RewriteKey["GET "+pattern]; exist {
				return errors.New("Can't register router:\"GET " + pattern + "\",it had been registed in controller:" + old)
			}
			actionName := tb.(string)
			r := createRewriteRouter("GET", pattern, controllerName, actionName)
			if r == nil {
				return errors.New("Can't register router:" + pattern)
			}
			srt.RewriteMap["GET"] = append(srt.RewriteMap["GET"], r)
			srt.RewriteKey["GET "+pattern] = controllerName
		case map[string]string:
			for method, action := range tb.(map[string]string) {
				method = strings.ToUpper(method)
				if old, exist := srt.RewriteKey[method+" "+pattern]; exist {
					return errors.New("Can't register router:\"" + method + " " + pattern + "\",it had been registed in controller:" + old)
				}
				r := createRewriteRouter(method, pattern, controllerName, action)
				if r == nil {
					return errors.New("Can't register router:\"" + method + " " + pattern + "\"")
				}
				srt.RewriteMap[method] = append(srt.RewriteMap[method], r)
				srt.RewriteKey[method+" "+pattern] = controllerName
			}
		}
	}
	return nil
}

// create rewrite router
func createRewriteRouter(method, pattern , controller, action string) *RewriteRouter {
	if pattern == "" || action == "" {
		return nil
	}
	urlParts := strings.Split(strings.Trim(pattern, "/"), "/")

	router := &RewriteRouter{
		urlParts:   make(map[int]map[string]string),
		staticNum:  0,
		method:     method,
		controller: controller,
		action:     action,
	}

	for idx, part := range urlParts {
		if part[0:1] == ":" {
			router.urlParts[idx] = map[string]string{
				"name": part[1:],
				"type": "var",
			}
		} else if part[0:1] == "*" { //正则?
			router.urlParts[idx] = map[string]string{
				"name": part,
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


func (srt *SpiderRouter) MatchRewrite(url, method string) (string, string, map[string]string, error) {
	if _, ok := srt.RewriteMap[method]; ok == false {
		return "", "", nil, errors.New("No match")
	}

	paths := strings.Split(strings.Trim(url, "/"), "/")
	for _, router := range srt.RewriteMap[method] {
		if match_param, match := srt.matchRouter(router, paths); !match {
			continue
		} else {
			return router.controller, router.action, match_param, nil
		}
	}

	return "", "", nil, errors.New("No match")
}

func (this *SpiderRouter) matchRouter(router *RewriteRouter, paths []string) (map[string]string, bool) {
	var match_param map[string]string
	var cnt int
	for idx, part := range paths {
		if _, exist := router.urlParts[idx]; !exist {
			return nil, false
		}
		if router.urlParts[idx]["type"] == "" {
			if router.urlParts[idx]["name"] == "*" {
				if idx < len(paths) {
					if match_param == nil {
						match_param = make(map[string]string)
					}

					param := paths[idx:]
					param_num := len(param) / 2
					for i := 0; i < param_num; i++ {
						match_param[param[i*2]] = param[i*2+1]
					}
				}
				break
			}

			if router.urlParts[idx]["name"] != part {
				return nil, false
			}
			cnt++
		} else if router.urlParts[idx]["type"] == "var" {
			if match_param == nil {
				match_param = make(map[string]string)
			}
			match_param[router.urlParts[idx]["name"]] = part
		}
	}

	if cnt != router.staticNum {
		return nil, false
	}
	return match_param, true
}

// Get controller and action name from
// request parame "m"
// eg: demo.index,return demo index
func (this *SpiderRouter) ParseMethod(method string) (controller_name string, action_name string) {
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
func (this *SpiderRouter) NewController(controllerName string) (reflect.Value, error) {
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