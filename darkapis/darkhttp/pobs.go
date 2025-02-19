package darkhttp

import (
	"net/http"

	"github.com/darklab8/fl-darkstat/darkapis/darkhttp/apiutils"
	"github.com/darklab8/fl-darkstat/darkcore/web"
	"github.com/darklab8/fl-darkstat/darkcore/web/registry"
	_ "github.com/darklab8/fl-darkstat/darkstat/configs_export"
)

// ShowAccount godoc
// @Summary      Getting list of Player Owned Bases
// @Description  in difference to Disco API, it is enriched with Nicknames/Infocard Names,Region names
// @Description  Sector coordinates, and extra information written in Infocard (totally reflecting Darkstat itself)
// @Tags         bases
// @Accept       json
// @Produce      json
// @Success      200  {array}  	configs_export.PoB
// @Router       /api/pobs [post]
func GetPoBs(webapp *web.Web, api *Api) *registry.Endpoint {
	return &registry.Endpoint{
		Url: "" + ApiRoute + "/pobs",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			if webapp.AppDataMutex != nil {
				webapp.AppDataMutex.Lock()
				defer webapp.AppDataMutex.Unlock()
			}

			apiutils.ReturnJson(&w, api.app_data.Configs.PoBs)
		},
	}
}

// ShowAccount godoc
// @Summary      PoB Goods
// @Tags         bases
// @Accept       json
// @Produce      json
// @Success      200  {array}  	configs_export.PoBGood
// @Router       /api/pob_goods [post]
func GetPobGoods(webapp *web.Web, api *Api) *registry.Endpoint {
	return &registry.Endpoint{
		Url: "" + ApiRoute + "/pob_goods",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			if webapp.AppDataMutex != nil {
				webapp.AppDataMutex.Lock()
				defer webapp.AppDataMutex.Unlock()
			}

			apiutils.ReturnJson(&w, api.app_data.Configs.PoBGoods)
		},
	}
}
