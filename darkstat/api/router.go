package api

import (
	"github.com/darklab8/fl-darkcore/darkcore/web"
	"github.com/darklab8/fl-darkcore/darkcore/web/registry"
	"github.com/darklab8/fl-darkstat/darkstat/router"
	httpSwagger "github.com/swaggo/http-swagger"
)

const ApiRoute = "/api"

type Api struct {
	app_data *router.AppData
}

func RegisterApiRoutes(w *web.Web, app_data *router.AppData) *web.Web {
	api := &Api{
		app_data: app_data,
	}
	api_routes := registry.NewRegister()
	api_routes.Register(NewEndpointPoBs(w, api))
	api_routes.Register(NewEndpointPoBGoods(w, api))
	api_routes.Register(NewEndpointHashes(w, api))

	w.GetMux().Handle("GET /swagger/", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))

	api_routes.Foreach(func(e *registry.Endpoint) {
		w.GetMux().HandleFunc(string(e.Url), e.Handler)
	})
	return w
}
