package editor

import (
	"github.com/lunny/tango"

	"github.com/Felamande/kiriadmin/middlewares/auth"
	"github.com/Felamande/kiriadmin/middlewares/xsrf"
	"github.com/tango-contrib/binding"

	"github.com/Felamande/kiriadmin/models"

	"github.com/Felamande/lotdb/routers/base"
)

// var _ xsrf.XsrfChecker = &PreviewRouter{}

type PreviewRouter struct {
	binding.Binder
	tango.Ctx
	auth.Auther
	xsrf.Checker

	JSON map[string]interface{}
	errs binding.Errors
	md   *models.Markdown
}

func (r *PreviewRouter) Before() {
	r.md = new(models.Markdown)
	if errs := r.Json(r.md); errs.Len() != 0 {
		r.errs = errs
	}

}

func (r *PreviewRouter) Post() {
	r.JSON = make(map[string]interface{})
	if r.errs.Len() != 0 {
		r.JSON["err"] = r.errs[0].Error()
		r.ServeJson(r.JSON)
		return
	}
	html := r.md.Convert()
	r.Write(html)
}
func (r *PreviewRouter) GetTokenFromReq() string {
	return r.md.Xsrf

}

func (r *PreviewRouter) DisposableToken() bool {
	return false
}

func (r *PreviewRouter) HandleXsrfErr(err error) {
	r.Logger.Info(err)
	r.errs.Add([]string{"xsrf"}, "xsrf", err.Error())
}

type EditorHome struct {
	base.BaseTplRouter
	auth.Auther
	xsrf.Checker
}

func (r *EditorHome) Get() {
	r.Data["title"] = "编辑器"
	r.Data["xsrf"] = r.CreateXsrfHTML()
	r.Tpl = "editor.html"
	r.Render(r.Tpl, r.Data)
}
