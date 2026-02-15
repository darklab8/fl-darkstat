package darkhttp

import (
	"net/http"

	"github.com/darklab8/fl-darkstat/darkapis/darkhttp/apiutils"
	"github.com/darklab8/fl-darkstat/darkcore/core_types"
	"github.com/darklab8/fl-darkstat/darkcore/web"
	"github.com/darklab8/fl-darkstat/darkcore/web/registry"
	"github.com/darklab8/fl-darkstat/darkstat/configs_export"
	_ "github.com/darklab8/fl-darkstat/darkstat/configs_export"
)

var GetPobsUrl = core_types.Url("" + ApiRoute + "/pobs")

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
		Url: GetPobsUrl,
		Handler: func(w http.ResponseWriter, r *http.Request) {
			if webapp.AppDataMutex != nil {
				webapp.AppDataMutex.RLock()
				defer webapp.AppDataMutex.RUnlock()
			}

			apiutils.ReturnJson(&w, api.app_data.Configs.PoBs)
		},
	}
}
func (c *HttpClient) GetPobs() ([]*configs_export.PoB, error) {
	return make_request[EmptyInput, []*configs_export.PoB](c, GetPobsUrl, EmptyInput{})
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
				webapp.AppDataMutex.RLock()
				defer webapp.AppDataMutex.RUnlock()
			}

			apiutils.ReturnJson(&w, api.app_data.Configs.PoBGoods)
		},
	}
}
func (c *HttpClient) GetPoBGoods() ([]*configs_export.PoBGood, error) {
	return make_request[EmptyInput, []*configs_export.PoBGood](c, ""+ApiRoute+"/pob_goods", EmptyInput{})
}
