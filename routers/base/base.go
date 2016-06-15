package base

import (
	"github.com/lunny/tango"
	"github.com/tango-contrib/renders"
)

type BaseRouter struct {
	tango.Ctx
}

type BaseJSONRouter struct {
	tango.Json
	tango.Ctx
	JSON map[string]interface{}
}

func (r *BaseJSONRouter) ReturnJSON() {
	r.ServeJson(r.JSON)
}

func (r *BaseJSONRouter) Before() {
	r.JSON = make(map[string]interface{})
}

func (r *BaseJSONRouter) After() {

}

type BaseTplRouter struct {
	tango.Ctx

	renders.Renderer
	Tpl  string
	Data renders.T
}

func (r *BaseTplRouter) Before() {
	r.Data = make(renders.T)
}

func (r *BaseTplRouter) After() {

}
