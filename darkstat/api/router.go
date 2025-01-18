package api

import (
	"encoding/json"
	"net/http"

	"github.com/darklab8/fl-darkcore/darkcore/web"
	"github.com/darklab8/fl-darkcore/darkcore/web/registry"
	"github.com/darklab8/fl-darkstat/darkstat/router"
	"github.com/darklab8/fl-darkstat/darkstat/settings/logus"
	httpSwagger "github.com/swaggo/http-swagger"
)

const ApiRoute = "/api"

type Api struct {
	app_data *router.AppData
}

func JsonResponseHeader(w *http.ResponseWriter) {
	(*w).Header().Set("Content-Type", "application/json")
}

func ReturnJson(w *http.ResponseWriter, data any) {
	(*w).Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(*w).Encode(data)
	if logus.Log.CheckError(err, "should be marshable") {
		json.NewEncoder(*w).Encode(struct {
			Error string
		}{
			Error: "not marshable for some reason",
		})
		(*w).WriteHeader(http.StatusInternalServerError)
	}
}

func RegisterApiRoutes(w *web.Web, app_data *router.AppData) *web.Web {
	api := &Api{
		app_data: app_data,
	}
	api_routes := registry.NewRegister()
	api_routes.Register(NewEndpointPoBs(w, api))
	api_routes.Register(NewEndpointPoBGoods(w, api))
	api_routes.Register(NewEndpointHashes(w, api))
	api_routes.Register(PostGraphPaths(w, api))
	api_routes.Register(NewEndpointBases(w, api))
	api_routes.Register(NewEndpointOreFields(w, api))

	w.GetMux().Handle("GET /swagger/", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))

	api_routes.Foreach(func(e *registry.Endpoint) {
		w.GetMux().HandleFunc(string(e.Url), e.Handler)
	})
	return w
}
