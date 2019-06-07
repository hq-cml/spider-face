package core

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

func (fc *SpeedyController) GetAllRouters() []RouteEntry {
	return fc.routers
}

func (fc *SpeedyController) GetRoundTrip() Roundtrip {
	return &fc.spdrt
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

func (fc *SpeedyController) SpeedyGetAction() {
	rp := fc.GetRoundTrip()
	actionFunc := fc.funcMapGet[rp.UrlPath()] //此处必然能存在，因为前面经过路径路由分析
	actionFunc(rp)
}

func (fc *SpeedyController) SpeedyPostAction() {
	rp := fc.GetRoundTrip()
	actionFunc := fc.funcMapPost[rp.UrlPath()]
	actionFunc(rp)
}

func (fc *SpeedyController) SpeedyPutAction() {
	rp := fc.GetRoundTrip()
	actionFunc := fc.funcMapPut[rp.UrlPath()]
	actionFunc(rp)
}

func (fc *SpeedyController) SpeedyDeleteAction() {
	rp := fc.GetRoundTrip()
	actionFunc := fc.funcMapDelete[rp.UrlPath()]
	actionFunc(rp)
}
