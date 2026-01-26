package darkhttp

import (
	"net/http"

	"github.com/darklab8/fl-darkstat/darkapis/darkgrpc_deprecated"
	"github.com/darklab8/fl-darkstat/darkcore/web"
	"github.com/darklab8/fl-darkstat/darkcore/web/registry"
	"github.com/darklab8/fl-darkstat/darkstat/appdata"
	"github.com/darklab8/fl-darkstat/darkstat/settings/logus"
	httpSwagger "github.com/swaggo/http-swagger"
)

var (
	Log = logus.Log.WithScope("darkhttp")
)

const ApiRoute = "/api"

type Api struct {
	app_data *appdata.AppData
}

func (a *Api) GetAppData() *appdata.AppData {
	return a.app_data
}

func JsonResponseHeader(w *http.ResponseWriter) {
	(*w).Header().Set("Content-Type", "application/json")
}

func RegisterApiRoutes(w *web.Web, app_data *appdata.AppData) *web.Web {
	api := &Api{
		app_data: app_data,
	}
	api_routes := registry.NewRegister()
	api_routes.Register(GetPoBs(w, api))
	api_routes.Register(GetPobGoods(w, api))
	api_routes.Register(darkgrpc_deprecated.GetPobGoodsDeprecated(w, api))
	api_routes.Register(GetHashes(w, api))
	api_routes.Register(PostGraphPaths(w, api))
	api_routes.Register(GetBases(w, api))
	api_routes.Register(GetOreFields(w, api))
	api_routes.Register(PostBaseMarketGoods(w, api))
	api_routes.Register(GetShips(w, api))
	api_routes.Register(PostShipMarketGoods(w, api))
	api_routes.Register(PostShipTechcompatibilities(w, api))
	api_routes.Register(GetCommodities(w, api))
	api_routes.Register(PostCommodityMarketGoods(w, api))
	api_routes.Register(GetTractors(w, api))
	api_routes.Register(PostTractorMarketGoods(w, api))
	api_routes.Register(GetFactions(w, api))
	api_routes.Register(GetAmmos(w, api))
	api_routes.Register(PostAmmoMarketGoods(w, api))
	api_routes.Register(PostAmmoTechcompatibilities(w, api))
	api_routes.Register(GetGuns(w, api))
	api_routes.Register(PostGunsMarketGoods(w, api))
	api_routes.Register(PostGunsTechcompatibilities(w, api))
	api_routes.Register(GetMissiles(w, api))
	api_routes.Register(PostMissilesMarketGoods(w, api))
	api_routes.Register(PostMissilesTechcompatibilities(w, api))
	api_routes.Register(GetMines(w, api))
	api_routes.Register(PostMinesMarketGoods(w, api))
	api_routes.Register(PostMinesTechcompatibilities(w, api))
	api_routes.Register(GetCMs(w, api))
	api_routes.Register(PostCMsMarketGoods(w, api))
	api_routes.Register(PostCMsTechcompatibilities(w, api))
	api_routes.Register(GetEngines(w, api))
	api_routes.Register(PostEnginesMarketGoods(w, api))
	api_routes.Register(PostEnginesTechcompatibilities(w, api))
	api_routes.Register(GetScanners(w, api))
	api_routes.Register(PostScannersMarketGoods(w, api))
	api_routes.Register(PostScannersTechcompatibilities(w, api))
	api_routes.Register(GetShields(w, api))
	api_routes.Register(PostShieldsMarketGoods(w, api))
	api_routes.Register(PostShieldsTechcompatibilities(w, api))
	api_routes.Register(GetThrusters(w, api))
	api_routes.Register(PostThrustersMarketGoods(w, api))
	api_routes.Register(PostThrustersTechcompatibilities(w, api))
	api_routes.Register(GetPoBBases(w, api))
	api_routes.Register(GetInfocards(w, api.app_data, api))
	api_routes.Register(GetSystems(w, api))

	w.GetMux().Handle("GET /swagger/", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))

	api_routes.Foreach(func(e *registry.Endpoint) {
		w.GetMux().HandleFunc(string(e.Url), e.Handler)
	})
	return w
}
