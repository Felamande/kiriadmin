package editor

import (
	"github.com/Felamande/kiriadmin/middlewares/auth"
	"github.com/Felamande/kiriadmin/models"
	"github.com/lunny/tango"
	"github.com/tango-contrib/binding"
)

type EvalMdRouter struct {
	binding.Binder
	tango.Ctx
	JSON map[string]interface{}
	auth.Auther
}

func (r *EvalMdRouter) Post() {

	md := models.Markdown{}
	r.JSON = make(map[string]interface{})
	if errs := r.Json(&md); errs.Len() != 0 {
		r.JSON["error"] = errs[0].Error()
		r.ServeJson(r.JSON)
		return
	}
	evaled := md.Eval()
	r.Write(evaled)

}
