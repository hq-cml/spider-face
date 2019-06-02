package core

/*
 * Controller接口，规定了Spider中必需的合法的行为
 */
type Controller interface {
	GetAllRouters() []ControllerRouter
	GetRoundTrip() Roundtrip
}

type ActionFunc func(rp Roundtrip)

type ControllerRouter struct {
	Method   string
	Location string
	Action   string
}

//系统自动注册一个默认的Controller，用于傻瓜式快捷注册
type FoolishController struct {
	spdrt           SpiderRoundtrip
	routers    		[]ControllerRouter

	funcMapGet  	map[string]ActionFunc
	funcMapPost  	map[string]ActionFunc
	funcMapPut  	map[string]ActionFunc
	funcMapDelete   map[string]ActionFunc
}

const FOOLISH_CONTROLLER_NAME = "Foolish"

func (fc *FoolishController) GetAllRouters() []ControllerRouter {
	return fc.routers
}

func (fc *FoolishController) GetRoundTrip() Roundtrip {
	return &fc.spdrt
}

func NewFoolishController() *FoolishController {
	return &FoolishController{
		routers:    	 []ControllerRouter{},
		funcMapGet:		 map[string]ActionFunc{},
		funcMapPost:	 map[string]ActionFunc{},
		funcMapPut: 	 map[string]ActionFunc{},
		funcMapDelete:	 map[string]ActionFunc{},
	}
}

func (fc *FoolishController) DefaultGetAction() {
	rp := fc.GetRoundTrip()
	actionFunc := fc.funcMapGet[rp.UrlPath()] //此处必然能存在，因为前面经过路径路由分析
	actionFunc(rp)
}

func (fc *FoolishController) DefaultPostAction() {
	rp := fc.GetRoundTrip()
	actionFunc := fc.funcMapPost[rp.UrlPath()]
	actionFunc(rp)
}

func (fc *FoolishController) DefaultPutAction() {
	rp := fc.GetRoundTrip()
	actionFunc := fc.funcMapPut[rp.UrlPath()]
	actionFunc(rp)
}

func (fc *FoolishController) DefaultDeleteAction() {
	rp := fc.GetRoundTrip()
	actionFunc := fc.funcMapDelete[rp.UrlPath()]
	actionFunc(rp)
}
