package core

import "net/http"

/*
 * Controller接口，规定了Spider中必需的合法的行为
 */
type Controller interface {
	GetAllRouters() []RouteEntry
	GetRoundTrip() Roundtrip
}

type ActionFunc func(rp Roundtrip)

type RouteEntry struct {
	Method   string
	Location string
	Action   string
}

//系统自动注册一个默认的Controller，用于傻瓜式快捷注册
type SpeedyController struct {
	spdrt           SpiderRoundtrip
	routers    		[]RouteEntry

	funcMapGet  	map[string]ActionFunc
	funcMapPost  	map[string]ActionFunc
	funcMapPut  	map[string]ActionFunc
	funcMapDelete   map[string]ActionFunc
}

const SPEEDY_CONTROLLER_NAME = "Speedy"

func (sc *SpeedyController) GetAllRouters() []RouteEntry {
	return sc.routers
}

func (sc *SpeedyController) GetRoundTrip() Roundtrip {
	return &sc.spdrt
}

func NewSpeedyController() *SpeedyController {
	return &SpeedyController{
		routers:    	 []RouteEntry{},
		funcMapGet:		 map[string]ActionFunc{},
		funcMapPost:	 map[string]ActionFunc{},
		funcMapPut: 	 map[string]ActionFunc{},
		funcMapDelete:	 map[string]ActionFunc{},
	}
}

//SpeedController的4个默认Action，所有的action函数都以这4个Action作为入口
func (sc *SpeedyController) SpeedyGetAction() {
	rp := sc.GetRoundTrip()
	actionFunc := sc.funcMapGet[rp.UrlPath()] //此处必然能存在，因为前面经过路径路由分析
	actionFunc(rp)
}

func (sc *SpeedyController) SpeedyPostAction() {
	rp := sc.GetRoundTrip()
	actionFunc := sc.funcMapPost[rp.UrlPath()]
	actionFunc(rp)
}

func (sc *SpeedyController) SpeedyPutAction() {
	rp := sc.GetRoundTrip()
	actionFunc := sc.funcMapPut[rp.UrlPath()]
	actionFunc(rp)
}

func (sc *SpeedyController) SpeedyDeleteAction() {
	rp := sc.GetRoundTrip()
	actionFunc := sc.funcMapDelete[rp.UrlPath()]
	actionFunc(rp)
}

//快捷注册，向SpeedController中注入路由规则
func (sc *SpeedyController) GET(location string , acFunc ActionFunc) {
	sc.routers = append(sc.routers, RouteEntry{
		Method: http.MethodGet, Location: location, Action:"SpeedyGetAction",
	})
	sc.funcMapGet[location] = acFunc
}

func (sc *SpeedyController) POST(location string , acFunc ActionFunc) {
	sc.routers = append(sc.routers, RouteEntry{
		Method: http.MethodPost, Location: location, Action:"SpeedyPostAction",
	})
	sc.funcMapPost[location] = acFunc
}

func (sc *SpeedyController) PUT(location string , acFunc ActionFunc) {
	sc.routers = append(sc.routers, RouteEntry{
		Method: http.MethodPut, Location: location, Action:"SpeedyPutAction",
	})
	sc.funcMapPut[location] = acFunc
}

func (sc *SpeedyController) DELETE(location string , acFunc ActionFunc) {
	sc.routers = append(sc.routers, RouteEntry{
		Method: http.MethodDelete, Location: location, Action:"SpeedyDeleteAction",
	})
	sc.funcMapDelete[location] = acFunc
}