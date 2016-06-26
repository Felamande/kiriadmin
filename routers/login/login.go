package login

import (
	"github.com/Felamande/kiriadmin/middlewares/auth"
	"github.com/Felamande/kiriadmin/models"
	"github.com/Felamande/kiriadmin/routers/base"
	"github.com/tango-contrib/captcha"
	"github.com/tango-contrib/renders"
	// "github.com/tango-contrib/xsrf"
)

var _ auth.Auther = &LoginRouter{}

type LoginRouter struct {
	base.BaseTplRouter
	captcha.Captcha
	auth.AuthAdmin

	JSON map[string]interface{}
}

func (r *LoginRouter) Get() {
	if r.IsLogin() {
		r.Redirect("/")
		return
	}
	r.Render("login.html", renders.T{
		"captcha": r.CreateHtml(),
		"title":   "登陆后台",
	})
}

func (r *LoginRouter) Post() {
	if r.JSON != nil {
		r.JSON = make(map[string]interface{})
	}
	r.Req().ParseForm()

	if !r.Verify() {
		r.JSON["err"] = "invalid captcha"
		r.ServeJson(r.JSON)
		r.Logger.Error(r.JSON["err"])
		return
	}

	name := r.Req().FormValue("name")
	pwd := r.Req().FormValue("pwd")
	u := models.GetUserByName(name)
	if u == nil {
		r.JSON["err"] = "no such user"
		r.ServeJson(r.JSON)
		r.Logger.Error(r.JSON["err"])
		return
	}
	if !u.Auth(pwd) {
		r.JSON["err"] = "invalid password"
		r.ServeJson(r.JSON)
		r.Logger.Error(r.JSON["err"])
		return
	}

	r.LoginWithToken("admin")
	r.Redirect("/")

}

func (r *LoginRouter) IsLogin() bool {
	if r.Token() == nil {
		return false
	}
	return "admin" == r.Token().(string)
}

func (r *LoginRouter) AskAuth() bool {
	return false
}

type LogoutRouter struct {
	base.BaseRouter
	auth.AuthAdmin
}

func (r *LogoutRouter) Get() {
	if r.IsLogin() {
		r.Logout()
	}
	r.Redirect("/login")

}

func (r *LogoutRouter) IsLogin() bool {
	return "admin" == r.Token().(string)
}
