package darkapi

import (
	"github.com/darklab8/fl-darkstat/darkcore/web"
	"github.com/darklab8/fl-darkstat/darkcore/web/registry"
	"github.com/darklab8/fl-darkstat/darkstat/configs_export"
)

type TractorrWithInfocard struct {
	configs_export.Tractor
	Infocard configs_export.Infocard
}

// ShowAccount godoc
// @Summary      Getting list of tractors
// @Tags         tractors
// @Accept       json
// @Produce      json
// @Success      200  {array}  	TractorrWithInfocard
// @Router       /api/tractors [get]
// @Param        filter_to_useful    query     string  false  "filter items only to useful, usually they are sold, or have goods, or craftable or findable in loot, or bases that are flight reachable from manhattan"  example("true")
func GetTractors(webapp *web.Web, api *Api) *registry.Endpoint {
	return &registry.Endpoint{
		Url:     "GET " + ApiRoute + "/tractors",
		Handler: GetItemsT(webapp, api.app_data, api.app_data.Configs.Tractors, api.app_data.Configs.FilterToUsefulTractors),
	}
}

// ShowAccount godoc
// @Summary      Getting list of Tractor Market Goods
// @Tags         tractors
// @Accept       json
// @Produce      json
// @Param request body []string true "Array of ship nicknames as input, for example [dsy_license_srp_28]"
// @Success      200  {array}  	MarketGoodResp
// @Router       /api/tractors/market_goods [post]
func PostTractorMarketGoods(webapp *web.Web, api *Api) *registry.Endpoint {
	return &registry.Endpoint{
		Url:     "POST " + ApiRoute + "/tractors/market_goods",
		Handler: PostItemsMarketGoodsT(webapp, api.app_data.Configs.Tractors),
	}
}
