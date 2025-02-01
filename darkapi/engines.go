package darkapi

import (
	"github.com/darklab8/fl-darkstat/darkcore/web"
	"github.com/darklab8/fl-darkstat/darkcore/web/registry"
)

// ShowAccount godoc
// @Summary      Getting list of Engines
// @Tags         engines
// @Accept       json
// @Produce      json
// @Success      200  {array}  	configs_export.Engine
// @Router       /api/engines [get]

func GetEngines(webapp *web.Web, api *Api) *registry.Endpoint {
	return &registry.Endpoint{
		Url:     "GET " + ApiRoute + "/engines",
		Handler: GetItemsT(webapp, api.app_data.Configs.Engines),
	}
}

// ShowAccount godoc
// @Summary      Getting list of Engine Market Goods
// @Tags         engines
// @Accept       json
// @Produce      json
// @Param request body []string true "Array of engines nicknames as input, for example [ge_kfr_engine_01_add]"
// @Success      200  {array}  	MarketGoodResp
// @Router       /api/engines/market_goods [post]
func PostEnginesMarketGoods(webapp *web.Web, api *Api) *registry.Endpoint {
	return &registry.Endpoint{
		Url:     "POST " + ApiRoute + "/engines/market_goods",
		Handler: PostItemsMarketGoodsT(webapp, api.app_data.Configs.Engines),
	}
}

// ShowAccount godoc
// @Summary      Getting list of Engine Tech compats
// @Tags         engines
// @Accept       json
// @Produce      json
// @Param request body []string true "Array of engines nicknames as input"
// @Success      200  {array}  	TechCompatResp
// @Router       /api/engines/tech_compats [post]
func PostEnginesTechcompatibilities(webapp *web.Web, api *Api) *registry.Endpoint {
	return &registry.Endpoint{
		Url:     "POST " + ApiRoute + "/engines/tech_compats",
		Handler: PostItemsTechCompatT(webapp, api.app_data.Configs.Engines),
	}
}
