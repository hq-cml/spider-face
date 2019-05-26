package core

import (
	"strings"
	"errors"
	"reflect"
)

type RouterManager struct {
	logger       SpiderLogger

	ParentMux    *SpiderHandlerMux
	RouterTable  map[string][]*RouterNode  //路由表，核心
	UniqKeyMap   map[string]string         //用于防重复校验
}

type RouterNode struct {
	UrlParts   []PathPartition
	NormalNum  int //整个path中normal的数量
	Method     string
	Controller string
	Action     string
}

//URL路径的一段，它可能是三种类型：
// 普通的字符，url中的普通的一段，没有什么特殊含义
// 动态参数，比如http://servername/user/birth/:year/:month/:day中的year,month,day
// Pathinfo参数，比如http://servername/user/birth/year/2019/month/12/day/23这种情况
type PathPartition struct {
	Type  string  //normal/variable/pathinfo
	Value string
}

func NewRouterManager(mux *SpiderHandlerMux, logger SpiderLogger) *RouterManager {
	return &RouterManager{
		logger:      logger,
		ParentMux:   mux,
		RouterTable: make(map[string][]*RouterNode),
		UniqKeyMap:  make(map[string]string),
	}
}

//注册路由
//将外部传入的Controller拆解成路由表的记录，注册注册进入Spider
func (rtm *RouterManager) RegisterRouter(controllerName string, controller SpiderController) error {
	routers := controller.GetRouter()
	if routers == nil {
		return nil
	}

	for _, router := range routers {
		method := router.Method
		pattern := router.Pattern
		action := router.Action

		method = strings.ToUpper(method)
		if old, exist := rtm.UniqKeyMap[method+" "+pattern]; exist {
			return errors.New("Can't register router" + method + " " + pattern + ". it had been registed in controller:" + old)
		}
		node, err := genRouterNode(method, pattern, controllerName, action)
		if err != nil {
			return err
		}
		if node == nil {
			return errors.New("Can't register router:" + method + " " + pattern)
		}
		rtm.RouterTable[method] = append(rtm.RouterTable[method], node)
		rtm.UniqKeyMap[method+" "+pattern] = controllerName
	}

	return nil
}

// create rewrite RouterNode
func genRouterNode(method, pattern, controller, action string) (*RouterNode, error) {
	if pattern == "" || action == "" {
		return nil, nil
	}
	urlParts := strings.Split(strings.Trim(pattern, "/"), "/")

	router := &RouterNode{
		UrlParts:   nil,
		NormalNum:  0,
		Method:     method,
		Controller: controller,
		Action:     action,
	}

	partitions := []PathPartition{}
	for idx, part := range urlParts {
		if part[0:1] == ":" {       //动态路径参数
			partitions = append(partitions, PathPartition{
				Value: part[1:],
				Type: "var",
			})
		} else if part == "*" { //Pathinfo的参数形式/yera/2019/month/5/day/10
			//pathinfo模式，*必须是最后一段
			if idx != len(urlParts) - 1 {
				return nil, errors.New("Pathinfo Error!")
			}
			partitions = append(partitions, PathPartition{
				Value: "*",
				Type: "pathinfo",
			})
			break
		} else {
			partitions = append(partitions, PathPartition{
				Value: part,
				Type: "normal",
			})
			router.NormalNum++
		}
	}
	router.UrlParts = partitions
	return router, nil
}

//根据URL和参数，找到对应处理的controller和action
func (rtm *RouterManager) AnalysePath(method, url string) (string, string, map[string]string, error) {
	if _, ok := rtm.RouterTable[method]; ok == false {
		rtm.logger.Errf("No support Method: %s", method)
		return "", "", nil, errors.New("No match")
	}

	paths := strings.Split(strings.Trim(url, "/"), "/")
	for _, rNode := range rtm.RouterTable[method] {
		if pathParam, match := rtm.matchRouter(rNode, paths); !match {
			continue
		} else {
			return rNode.Controller, rNode.Action, pathParam, nil
		}
	}

	return "", "", nil, errors.New("No match")
}

//判断一个给定的rNode和路径paths是否匹配
func (rtm *RouterManager) matchRouter(rNode *RouterNode, paths []string) (map[string]string, bool) {
	pathParam := map[string]string{}
	var cnt int
	for idx, part := range paths {
		if idx > len(rNode.UrlParts)-1 {
			return nil, false
		}
		if rNode.UrlParts[idx].Type == "pathinfo" {
			if idx < len(paths)-1 {
				param := paths[idx:]
				paramNum := len(param) / 2
				for i := 0; i < paramNum; i++ {
					pathParam[param[i*2]] = param[i*2+1]
				}
			}
			break
		} else if rNode.UrlParts[idx].Type == "var" {
			pathParam[rNode.UrlParts[idx].Value] = part
		} else {
			if rNode.UrlParts[idx].Value != part {
				return nil, false
			}
			cnt ++
		}
	}

	if cnt != rNode.NormalNum {
		return nil, false
	}
	return pathParam, true
}

// create new controller by controller name
func (rtm *RouterManager) NewController(controllerName string) (reflect.Value, error) {
	//register := GetRegister()

	//m_arr := make([]reflect.Value, 0)
	var newController reflect.Value

	//if register == nil {
	//	//http 500
	//	return controller_instance, errors.New("Server Error : Can't find \"Register\"")
	//}

	controller_type := rtm.ParentMux.GetController(controllerName)
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