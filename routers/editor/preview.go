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
	auth.AuthAdmin
	xsrf.Checker

	JSON map[string]interface{}
	errs binding.Errors
	md   models.Markdown
}

func (r *PreviewRouter) Before() {
	if errs := r.Json(&r.md); errs.Len() != 0 {
		r.errs = errs
	}

}

func (r *PreviewRouter) Post() {
	r.JSON = make(map[string]interface{})
	if r.errs.Len() != 0 {
		r.JSON["err"] = r.errs[0].Error()
		r.ServeJson(r.JSON)
		r.Logger.Info(r.errs)
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
	r.errs.Add([]string{"xsrf"}, "xsrf", err.Error())
}

type EditorHome struct {
	base.BaseTplRouter
	auth.AuthAdmin
	xsrf.Checker
}

func (r *EditorHome) Get() {
	r.Data["title"] = "编辑器"
	r.Data["xsrf"] = r.CreateXsrfHTML()
	r.Tpl = "editor.html"
	r.Render(r.Tpl, r.Data)
}

type CommitRouter struct {
	binding.Binder
	xsrf.Checker
	auth.AuthAdmin
	errs binding.Errors

	article models.Article
}

func (r *CommitRouter) Before() {
	r.errs = r.Json(&r.article)
}

func (r *CommitRouter) Post() {

}

const articleTmp = `title: {{.Title}}
date: {{.Date}} +0800
update: {{.Update}} +0800
author: {{.Author}}
{{if .Cover}}Cover: "{{.Cover}}"{{end}}
{{if .Tags}}
tags:
    {{range _,$tag := .Tags}}
	- {{$tag}}
	{{end}}
{{end}}
{{if .Preview}}preview:{{.Preview}}{{end}}

---
{{.Content}}
`
