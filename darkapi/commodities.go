package darkapi

import (
	"github.com/darklab8/fl-darkstat/darkcore/web"
	"github.com/darklab8/fl-darkstat/darkcore/web/registry"
)

// ShowAccount godoc
// @Summary      Getting list of Commodities
// @Tags         commodities
// @Accept       json
// @Produce      json
// @Success      200  {array}  	configs_export.Commodity
// @Router       /api/commodities [get]
// @Param        filter_to_useful    query     string  false  "filter items only to useful, usually they are sold, or have goods, or craftable or findable in loot, or bases that are flight reachable from manhattan"
func GetCommodities(webapp *web.Web, api *Api) *registry.Endpoint {
	return &registry.Endpoint{
		Url:     "GET " + ApiRoute + "/commodities",
		Handler: GetItemsT(webapp, api.app_data.Configs.Commodities, api.app_data.Configs.FilterToUsefulCommodities),
	}
}

// ShowAccount godoc
// @Summary      Getting list of Commodities Market Goods
// @Tags         commodities
// @Accept       json
// @Produce      json
// @Param request body []string true "Array of commodity nicknames as input, for example [commodity_military_salvage]"
// @Success      200  {array}  	MarketGoodResp
// @Router       /api/commodities/market_goods [post]
func PostCommodityMarketGoods(webapp *web.Web, api *Api) *registry.Endpoint {
	return &registry.Endpoint{
		Url:     "POST " + ApiRoute + "/commodities/market_goods",
		Handler: PostItemsMarketGoodsT(webapp, api.app_data.Configs.Commodities),
	}
}
