package core

import (
	"strings"
	"errors"
)

/*
 * 维护路由表,
 */
type RouterManager struct {
	logger       SpiderLogger

	RouterTable  map[string][]*RouteNode //路由表，核心 method=>list(node)
	UniqKeyMap   map[string]string       //用于防重复校验
}

//一条路由规则
type RouteNode struct {
	UrlParts   []PathPartition
	NormalNum  int             //整个path中normal段的数量
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

func NewRouterManager(logger SpiderLogger) *RouterManager {
	return &RouterManager{
		logger:      logger,
		RouterTable: make(map[string][]*RouteNode),
		UniqKeyMap:  make(map[string]string),
	}
}

//注册路由
//将外部传入的Controller拆解成路由表的记录，注册注册进入Spider
func (rtm *RouterManager) RegisterRouter(controllerName string, controller Controller) error {
	routers := controller.GetAllRouters()
	if routers == nil {
		return nil
	}

	for _, router := range routers {
		method := router.Method
		pattern := router.Location
		action := router.Action

		method = strings.ToUpper(method)
		if old, exist := rtm.UniqKeyMap[method + " " + pattern]; exist {
			return errors.New("Can't register router" + method + " " + pattern + ". it had been registed in controller:" + old)
		}
		node, err := genRouterNode(pattern, controllerName, action)
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

// create rewrite RouteNode
func genRouterNode(pattern, controller, action string) (*RouteNode, error) {
	if pattern == "" || action == "" {
		return nil, nil
	}

	var urlParts []string
	if pattern == "/" {
		urlParts = []string{"/"}
	} else  {
		urlParts = strings.Split(strings.Trim(pattern, "/"), "/")
	}

	router := &RouteNode{
		UrlParts:   nil,
		NormalNum:  0,
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
		} else if part == PATH_INFO_IDENTITY { //Pathinfo的参数形式/yera/2019/month/5/day/10
			//pathinfo模式，*必须是最后一段
			if idx != len(urlParts) - 1 {
				return nil, errors.New("Pathinfo Error!")
			}
			partitions = append(partitions, PathPartition{
				Value: PATH_INFO_IDENTITY,
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

	paths := []string{}
	if url != "/" {
		url = strings.Trim(url, "/")
		paths = strings.Split(url, "/")
	} else {
		paths = append(paths, "/")
	}
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
func (rtm *RouterManager) matchRouter(rNode *RouteNode, paths []string) (map[string]string, bool) {
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
