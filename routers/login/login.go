package login

import (
	"github.com/Felamande/kiriadmin/middlewares/auth"
	"github.com/Felamande/kiriadmin/models"
	"github.com/Felamande/kiriadmin/routers/base"
	"github.com/tango-contrib/captcha"
	"github.com/tango-contrib/renders"
	"github.com/tango-contrib/xsrf"
)

type LoginRouter struct {
	base.BaseTplRouter
	captcha.Captcha
	xsrf.Checker
	auth.Auther
}

type J map[string]interface{}

func (r *LoginRouter) Get() {
	r.Render("login.html", renders.T{
		"xsrfHTML": r.XsrfFormHtml(),
		"captcha":  r.CreateHtml(),
	})
}

func (r *LoginRouter) Post() {
	r.Req().ParseForm()

	if !r.Verify() {
		r.Logger.Error("invalid captcha")
		r.ServeJson(J{"error": "invalid captcha", "sucess": false})
		return
	}

	name := r.Req().FormValue("name")
	pwd := r.Req().FormValue("pwd")
	u := models.GetUserByName(name)
	if u == nil {
		r.Logger.Error("invalid user")
		r.ServeJson(J{"error": "no such user", "sucess": false})
		return
	}
	if !u.Auth(pwd) {
		r.Logger.Error("invalid pwd")
		r.ServeJson(J{"error": "invalid password", "sucess": false})
		return
	}

	r.LoginWithToken("admin")
	r.Redirect("/")

}

func (r *LoginRouter) CheckToken() bool {
	return "admin" == r.Token()
}

func (r *LoginRouter) AskAuth() bool {
	return false
}
