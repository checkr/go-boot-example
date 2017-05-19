package controller

import (
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/checkr/go-boot-example/core"
)

type Ping struct{}

func init() {
	core.RegisterController(func() core.Controller {
		return &Ping{}
	})
}

func (x *Ping) Routes(api *rest.Api) []*rest.Route {
	var routes []*rest.Route

	routes = append(routes,
		rest.Get("/ping", x.get),
		rest.Post("/ping", x.post),
	)

	return routes
}

func (x *Ping) get(w rest.ResponseWriter, r *rest.Request) {
	w.WriteJson("GET")
}

func (x *Ping) post(w rest.ResponseWriter, r *rest.Request) {
	w.WriteJson("POST")
}
