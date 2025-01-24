package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/darklab8/fl-darkcore/darkcore/web"
	"github.com/darklab8/fl-darkcore/darkcore/web/registry"
	_ "github.com/darklab8/fl-darkstat/darkstat/configs_export"
	"github.com/darklab8/fl-darkstat/darkstat/settings/logus"
)

// ShowAccount godoc
// @Summary      Getting list of Player Owned Bases
// @Description  in difference to Disco API, it is enriched with Nicknames/Infocard Names,Region names
// @Description  Sector coordinates, and extra information written in Infocard (totally reflecting Darkstat itself)
// @Tags         pobs
// @Accept       json
// @Produce      json
// @Success      200  {array}  	configs_export.PoB
// @Router       /api/pobs [get]
func GetPoBs(webapp *web.Web, api *Api) *registry.Endpoint {
	return &registry.Endpoint{
		Url: "GET " + ApiRoute + "/pobs",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			if webapp.AppDataMutex != nil {
				webapp.AppDataMutex.Lock()
				defer webapp.AppDataMutex.Unlock()
			}

			ReturnJson(&w, api.app_data.Configs.PoBs)
		},
	}
}

// ShowAccount godoc
// @Summary      PoB Goods
// @Tags         pobs
// @Accept       json
// @Produce      json
// @Success      200  {array}  	configs_export.PoBGood
// @Router       /api/pob_goods [get]
func GetPobGoods(webapp *web.Web, api *Api) *registry.Endpoint {
	return &registry.Endpoint{
		Url: "GET " + ApiRoute + "/pob_goods",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			if webapp.AppDataMutex != nil {
				webapp.AppDataMutex.Lock()
				defer webapp.AppDataMutex.Unlock()
			}

			data, err := json.Marshal(api.app_data.Configs.PoBGoods)
			logus.Log.CheckPanic(err, "should be marshable")
			fmt.Fprint(w, string(data))
		},
	}
}
