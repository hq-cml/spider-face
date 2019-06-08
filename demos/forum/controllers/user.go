package controllers

import (
	"github.com/hq-cml/spider-face/core"

	"github.com/hq-cml/spider-face/demos/forum/model"
	"fmt"
)

//创建标准的Controller，必须以Controller后缀结尾
//并实现core.Controller接口
type UserController struct {
	spdrp core.SpiderRoundtrip
	entries []core.RouteEntry
}
func (uc *UserController) GetAllRouters() []core.RouteEntry {
	return uc.entries
}
func (uc *UserController) GetRoundTrip() core.Roundtrip {
	return &uc.spdrp
}

//创建一个controller
func NewUserAction() *UserController{
	uc := &UserController{}
	return uc
}

//设置路由规则
func (uc *UserController) SetRouteEntries(entries []core.RouteEntry) {
	uc.entries = entries
}

//登陆页面
func (uc *UserController) LoginAction(rp core.Roundtrip) {
	rp.Display("user/login")
}

//登出
func (uc *UserController) LogoutAction(rp core.Roundtrip) {
	cookie := rp.GetCookie("_cookie")
	sess := model.Session{
		Uuid: cookie,
	}
	sess.DeleteByUUID()
	rp.Redirect("/index")
}

//注册页面
func (uc *UserController) SignupAction(rp core.Roundtrip) {
	rp.Display("user/signup")
}

//注册
func (uc *UserController) SignupAccountAction(rp core.Roundtrip) {
	user := model.User{
		Name: rp.Param("name"),
		Email: rp.Param("email"),
		Password: rp.Param("password"),
	}

	_, err := user.Create()
	if err != nil {
		msg := fmt.Sprintf("Create user Error: %v", err)
		rp.Redirect(fmt.Sprintf("/err?msg=%s", msg))
		return
	}

	rp.Redirect("/login")
}

//登陆的处理逻辑
func (uc *UserController) AuthenticateAction(rp core.Roundtrip) {
	email := rp.Param("email")
	pwd := rp.Param("password")

	user, err := model.GetUserByEmail(email)
	if err != nil {
		msg := fmt.Sprintf("Can't got user: %v", email)
		rp.Redirect(fmt.Sprintf("/err?msg=%s", msg))
		return
	}

	if user.Password == model.Encrypt(pwd) {
		//登陆成功
		sess, err := user.CreateSession()
		if err != nil {
			core.OutputErrorHtml(rp.GetResponse(), rp.GetRequest(), 500, nil)
		}
		rp.SetCookie("_cookie", sess.Uuid)
		rp.Redirect(fmt.Sprintf("/index"))
	} else {
		//登录失败
		rp.Redirect("/login")
	}
}

