package home

import (
	"github.com/Felamande/kiriadmin/middlewares/auth"
	"github.com/Felamande/kiriadmin/routers/base"
	"github.com/tango-contrib/renders"
)

type HomeRouter struct {
	base.BaseTplRouter
	auth.Auther
}

func (r *HomeRouter) Get() {
	if r.Data == nil {
		r.Data = make(renders.T)
	}
	r.Data["title"] = "管理员后台"
	r.Tpl = "admin.html"

	r.Render(r.Tpl, r.Data)
}
